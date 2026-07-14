package service

import (
	"context"
	"log"
	"strings"
	"time"

	"neocentral-go/auth-service/internal/domain"
	"neocentral-go/auth-service/internal/dto"
	"neocentral-go/auth-service/internal/event"
	"neocentral-go/auth-service/internal/repository"
	"neocentral-go/pkg/apperror"
	"neocentral-go/pkg/auth"

	"github.com/redis/go-redis/v9"
)

// AuthService contains all authentication-related business logic.
type AuthService struct {
	userRepo   repository.UserRepository
	jwtCfg     auth.JWTConfig
	events     *event.UserEventPublisher
	redis      *redis.Client
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	repo repository.UserRepository,
	jwtCfg auth.JWTConfig,
	events *event.UserEventPublisher,
	rdb *redis.Client,
) *AuthService {
	return &AuthService{
		userRepo: repo,
		jwtCfg:   jwtCfg,
		events:   events,
		redis:    rdb,
	}
}

// ── Login ────────────────────────────────────────────────────────

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if user == nil || user.Password == nil {
		return nil, apperror.Unauthorized("Invalid credentials")
	}

	if !user.IsVerified {
		return nil, apperror.New(403, "Akun belum diaktivasi. Silakan aktivasi akun terlebih dahulu.")
	}

	if !auth.CheckPassword(req.Password, *user.Password) {
		return nil, apperror.Unauthorized("Invalid credentials")
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID, derefStr(user.Email), s.jwtCfg)
	if err != nil {
		return nil, apperror.InternalWrap("failed to generate access token", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, derefStr(user.Email), s.jwtCfg)
	if err != nil {
		return nil, apperror.InternalWrap("failed to generate refresh token", err)
	}

	// Store hashed refresh token
	refreshHash := auth.HashToken(refreshToken)
	_ = s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"refresh_token": refreshHash,
	})

	// Publish login event (best-effort)
	if s.events != nil {
		go func() {
			if err := s.events.PublishUserLogin(user.ID, derefStr(user.Email)); err != nil {
				log.Printf("⚠️  Failed to publish login event: %v", err)
			}
		}()
	}

	return &dto.LoginResponse{
		User:         buildProfileResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ── Refresh Tokens ───────────────────────────────────────────────

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := auth.VerifyRefreshToken(refreshToken, s.jwtCfg.RefreshSecret)
	if err != nil {
		return nil, apperror.Unauthorized("Invalid refresh token")
	}

	user, err := s.userRepo.FindByID(ctx, claims.Subject)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if user == nil || user.RefreshToken == nil {
		return nil, apperror.Unauthorized("Invalid refresh token")
	}

	// Verify stored hash matches
	if !auth.CheckToken(refreshToken, *user.RefreshToken) {
		return nil, apperror.Unauthorized("Invalid refresh token")
	}

	// Rotate tokens
	newAccess, err := auth.GenerateAccessToken(user.ID, derefStr(user.Email), s.jwtCfg)
	if err != nil {
		return nil, apperror.InternalWrap("failed to generate access token", err)
	}

	newRefresh, err := auth.GenerateRefreshToken(user.ID, derefStr(user.Email), s.jwtCfg)
	if err != nil {
		return nil, apperror.InternalWrap("failed to generate refresh token", err)
	}

	newHash := auth.HashToken(newRefresh)
	_ = s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"refresh_token": newHash,
	})

	return &dto.TokenResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

// ── Logout ───────────────────────────────────────────────────────

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"refresh_token": nil,
	})
}

// ── Change Password ──────────────────────────────────────────────

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return apperror.InternalWrap("database error", err)
	}
	if user == nil || user.Password == nil {
		return apperror.NotFound("User not found")
	}

	if !auth.CheckPassword(req.CurrentPassword, *user.Password) {
		return apperror.BadRequest("Current password is incorrect")
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.InternalWrap("failed to hash password", err)
	}

	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"password":      hash,
		"refresh_token": nil, // force re-login on other sessions
	})
}

// ── Password Reset Request ───────────────────────────────────────

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, strings.ToLower(email))
	if err != nil || user == nil || user.Email == nil {
		// Return nil to avoid user enumeration
		return nil
	}

	token, err := auth.GeneratePurposeToken(user.ID, "pwdreset", 15*time.Minute, s.jwtCfg)
	if err != nil {
		return apperror.InternalWrap("failed to generate reset token", err)
	}

	// Store token flag in Redis (15 min TTL)
	key := "pwdreset:" + user.ID
	s.redis.Set(ctx, key, "1", 15*time.Minute)

	// TODO: Send email with reset link containing `token`
	log.Printf("🔗 Password reset token for %s: %s", derefStr(user.Email), token)

	return nil
}

// ── Reset Password with Token ────────────────────────────────────

func (s *AuthService) ResetPasswordWithToken(ctx context.Context, tokenStr, newPassword string) error {
	claims, err := auth.VerifyAccessToken(tokenStr, s.jwtCfg.Secret)
	if err != nil || claims.Purpose != "pwdreset" {
		return apperror.BadRequest("Invalid or expired token")
	}

	// Check Redis flag
	key := "pwdreset:" + claims.Subject
	exists, _ := s.redis.Get(ctx, key).Result()
	if exists == "" {
		return apperror.BadRequest("Token expired or already used")
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return apperror.InternalWrap("failed to hash password", err)
	}

	err = s.userRepo.UpdateFields(ctx, claims.Subject, map[string]interface{}{
		"password":      hash,
		"refresh_token": nil,
	})
	if err != nil {
		return apperror.InternalWrap("failed to update password", err)
	}

	// Consume token
	s.redis.Del(ctx, key)
	return nil
}

// ── Account Verification ─────────────────────────────────────────

func (s *AuthService) VerifyAccount(ctx context.Context, tokenStr string) error {
	claims, err := auth.VerifyAccessToken(tokenStr, s.jwtCfg.Secret)
	if err != nil || claims.Purpose != "verify" {
		return apperror.BadRequest("Invalid or expired token")
	}

	key := "verify:" + claims.Subject
	exists, _ := s.redis.Get(ctx, key).Result()
	if exists == "" {
		return apperror.BadRequest("Token expired or already used")
	}

	err = s.userRepo.UpdateFields(ctx, claims.Subject, map[string]interface{}{
		"is_verified": true,
	})
	if err != nil {
		return apperror.InternalWrap("failed to verify account", err)
	}

	s.redis.Del(ctx, key)
	return nil
}

// ── Request Account Verification ─────────────────────────────────

func (s *AuthService) RequestAccountVerification(ctx context.Context, email string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if user == nil || user.Email == nil {
		return map[string]interface{}{
			"found":   false,
			"message": "Email tidak terdaftar. Silakan hubungi admin untuk aktivasi akun.",
		}, nil
	}

	if user.IsVerified {
		return map[string]interface{}{
			"found":           true,
			"alreadyVerified": true,
			"message":         "Akun sudah terverifikasi.",
		}, nil
	}

	// Generate verification token
	token, err := auth.GeneratePurposeToken(user.ID, "verify", 24*time.Hour, s.jwtCfg)
	if err != nil {
		return nil, apperror.InternalWrap("failed to generate verification token", err)
	}

	// Store in Redis (24h TTL)
	s.redis.Set(ctx, "verify:"+user.ID, "1", 24*time.Hour)

	// TODO: Send verification email
	log.Printf("🔗 Verification token for %s: %s", derefStr(user.Email), token)

	return map[string]interface{}{
		"found":           true,
		"alreadyVerified": false,
		"sent":            true,
	}, nil
}

// ── Verify Access Token (for middleware / gRPC) ──────────────────

func (s *AuthService) VerifyAccessToken(tokenStr string) (*auth.Claims, error) {
	return auth.VerifyAccessToken(tokenStr, s.jwtCfg.Secret)
}

// ── Helpers ──────────────────────────────────────────────────────

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func buildProfileResponse(user *domain.User) dto.UserProfileResponse {
	resp := dto.UserProfileResponse{
		ID:             user.ID,
		FullName:       user.FullName,
		Email:          user.Email,
		IdentityNumber: user.IdentityNumber,
		IdentityType:   string(user.IdentityType),
		PhoneNumber:    user.PhoneNumber,
		IsVerified:     user.IsVerified,
		AvatarURL:      user.AvatarURL,
		Roles:          make([]dto.RoleResponse, 0, len(user.UserHasRoles)),
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      user.UpdatedAt.Format(time.RFC3339),
	}

	for _, uhr := range user.UserHasRoles {
		resp.Roles = append(resp.Roles, dto.RoleResponse{
			ID:     uhr.RoleID,
			Name:   uhr.Role.Name,
			Status: string(uhr.Status),
		})
	}

	if user.Student != nil {
		resp.Student = &dto.StudentResponse{
			ID:             user.Student.ID,
			EnrollmentYear: user.Student.EnrollmentYear,
			SKSCompleted:   user.Student.SKSCompleted,
			Status:         string(user.Student.Status),
		}
	}

	if user.Lecturer != nil {
		resp.Lecturer = &dto.LecturerResponse{
			ID:             user.Lecturer.ID,
			ScienceGroupID: user.Lecturer.ScienceGroupID,
		}
	}

	return resp
}
