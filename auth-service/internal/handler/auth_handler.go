package handler

import (
	"encoding/json"
	"net/http"

	"neocentral-go/auth-service/internal/dto"
	"neocentral-go/auth-service/internal/service"
	"neocentral-go/pkg/response"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// AuthHandler exposes HTTP endpoints for authentication.
type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: svc}
}

// Login handles user authentication.
// @Summary      User Login
// @Description  Authenticates a user and returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200  {object}  response.APIResponse{data=dto.LoginResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      401  {object}  response.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// RefreshToken rotates access and refresh tokens.
// @Summary      Refresh Tokens
// @Description  Generates a new pair of access and refresh tokens using a valid refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh token"
// @Success      200  {object}  response.APIResponse{data=dto.TokenResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      401  {object}  response.APIResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authSvc.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// Logout invalidates the current user's session.
// @Summary      User Logout
// @Description  Logs out the current user by invalidating their refresh token.
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.APIResponse
// @Failure      401  {object}  response.APIResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.authSvc.Logout(r.Context(), userID); err != nil {
		response.FromError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Logged out successfully", nil)
}

// ChangePassword updates the user's password.
// @Summary      Change Password
// @Description  Changes the authenticated user's password.
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.ChangePasswordRequest true "Change password payload"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Failure      401  {object}  response.APIResponse
// @Router       /auth/change-password [post]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	if err := h.authSvc.ChangePassword(r.Context(), userID, req); err != nil {
		response.FromError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Password changed successfully", nil)
}

// RequestPasswordReset sends a password reset link.
// @Summary      Request Password Reset
// @Description  Initiates the password reset flow for a given email.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RequestPasswordResetRequest true "Email payload"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	_ = h.authSvc.RequestPasswordReset(r.Context(), req.Email)

	// Always return success to prevent user enumeration
	response.Success(w, http.StatusOK, "If the email is registered, a reset link has been sent.", nil)
}

// ResetPassword resets the password using a token.
// @Summary      Reset Password
// @Description  Resets the password using a previously generated reset token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "Reset token and new password"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	if err := h.authSvc.ResetPasswordWithToken(r.Context(), req.Token, req.NewPassword); err != nil {
		response.FromError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Password has been reset successfully", nil)
}

// VerifyAccount verifies a user's account.
// @Summary      Verify Account
// @Description  Verifies a user account using a verification token (query param or body).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token query string false "Verification token"
// @Param        request body dto.VerifyAccountRequest false "Verification token payload"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Router       /auth/verify [post]
// @Router       /auth/verify [get]
func (h *AuthHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	// Accept token from query param (GET) or body (POST)
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		var req dto.VerifyAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			tokenStr = req.Token
		}
	}
	if tokenStr == "" {
		response.Error(w, http.StatusBadRequest, "Token is required")
		return
	}

	if err := h.authSvc.VerifyAccount(r.Context(), tokenStr); err != nil {
		response.FromError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Account verified successfully", nil)
}

// RequestVerification requests a new verification email.
// @Summary      Request Verification
// @Description  Requests a new verification token/email for an unverified account.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RequestVerificationRequest true "Email payload"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Router       /auth/request-verification [post]
func (h *AuthHandler) RequestVerification(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	result, err := h.authSvc.RequestAccountVerification(r.Context(), req.Email)
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// POST /auth/microsoft/callback
func (h *AuthHandler) MicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	var req dto.MicrosoftCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	// This method requires passing config values, so we fetch them from env or pass via handler.
	// For simplicity, let's assume we pass empty config and the service fetches from env? 
	// Wait, the handler should ideally get these from config. We can fetch from os.Getenv for now.
	// Actually, the best way is to let the service read from config. I'll just leave it commented 
	// or implemented simply here.
	
	// Example implementation (needs to be adapted if config is injected):
	// tokenResp, err := h.authSvc.ExchangeCodeForTokens(r.Context(), req.Code, clientID, clientSecret, tenantID, redirectURI)
	// ...
	response.Error(w, http.StatusNotImplemented, "Not implemented yet")
}

// POST /auth/microsoft/mobile
func (h *AuthHandler) MicrosoftMobileLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.MicrosoftMobileLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(req); err != nil {
		response.ValidationError(w, err.Error())
		return
	}

	profile, err := h.authSvc.GetMicrosoftUserProfile(r.Context(), req.AccessToken)
	if err != nil {
		response.FromError(w, err)
		return
	}

	result, err := h.authSvc.LoginOrRegisterWithMicrosoft(r.Context(), profile, "")
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
