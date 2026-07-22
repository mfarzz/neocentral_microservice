package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"neocentral-go/auth-service/internal/dto"
	"neocentral-go/pkg/apperror"
	"neocentral-go/pkg/auth"
)

const (
	graphAPIBase      = "https://graph.microsoft.com/v1.0"
	microsoftTokenURL = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
)

type MicrosoftTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type MicrosoftUserProfile struct {
	ID                string `json:"id"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	DisplayName       string `json:"displayName"`
}

// ExchangeCodeForTokens exchanges the OAuth2 auth code for tokens from Microsoft.
func (s *AuthService) ExchangeCodeForTokens(ctx context.Context, code, clientID, clientSecret, tenantID, redirectURI string) (*MicrosoftTokenResponse, error) {
	tokenURL := fmt.Sprintf(microsoftTokenURL, tenantID)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("scope", "user.read openid profile email offline_access")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("microsoft token exchange failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var tokenResp MicrosoftTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetMicrosoftUserProfile fetches the user's profile from Microsoft Graph API.
func (s *AuthService) GetMicrosoftUserProfile(ctx context.Context, accessToken string) (*MicrosoftUserProfile, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", graphAPIBase+"/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("microsoft graph request failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var profile MicrosoftUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// LoginOrRegisterWithMicrosoft handles the NeoCentral login logic given a Microsoft profile.
func (s *AuthService) LoginOrRegisterWithMicrosoft(ctx context.Context, profile *MicrosoftUserProfile, msRefreshToken string) (*dto.LoginResponse, error) {
	email := profile.Mail
	if email == "" {
		email = profile.UserPrincipalName
	}
	if email == "" {
		return nil, apperror.Unauthorized("Email not found in Microsoft account")
	}

	user, err := s.userRepo.FindByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}

	if user == nil {
		return nil, apperror.NotFound("Akun belum terdaftar. Silakan hubungi admin.")
	}

	if !user.IsVerified {
		return nil, apperror.Forbidden("Akun belum diaktivasi. Silakan aktivasi akun terlebih dahulu.")
	}

	// Update OAuth info
	err = s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"oauth_provider":      "microsoft",
		"oauth_id":            profile.ID,
		"oauth_refresh_token": msRefreshToken,
	})
	if err != nil {
		return nil, apperror.InternalWrap("failed to update user oauth info", err)
	}

	// Re-fetch to get updated data and relations
	user, _ = s.userRepo.FindByID(ctx, user.ID)

	var roles []string
	for _, uhr := range user.UserHasRoles {
		roles = append(roles, uhr.Role.Name)
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, derefStr(user.Email), roles, s.jwtCfg)
	if err != nil {
		return nil, apperror.InternalWrap("failed to generate access token", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, derefStr(user.Email), s.jwtCfg)
	if err != nil {
		return nil, apperror.InternalWrap("failed to generate refresh token", err)
	}

	refreshHash, _ := auth.HashPassword(refreshToken)
	_ = s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"refresh_token": refreshHash,
	})

	// Publish login event
	if s.events != nil {
		go s.events.PublishUserLogin(user.ID, derefStr(user.Email))
	}

	return &dto.LoginResponse{
		User:         buildProfileResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
