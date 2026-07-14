package dto

// ── Auth Requests ────────────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

type VerifyAccountRequest struct {
	Token string `json:"token" validate:"required"`
}

type RequestVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type MicrosoftCallbackRequest struct {
	Code     string `json:"code" validate:"required"`
	IsMobile bool   `json:"isMobile"`
}

type MicrosoftMobileLoginRequest struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

// ── Auth Responses ───────────────────────────────────────────────

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	User         UserProfileResponse `json:"user"`
	AccessToken  string              `json:"accessToken"`
	RefreshToken string              `json:"refreshToken"`
}

// ── Profile ──────────────────────────────────────────────────────

type UserProfileResponse struct {
	ID             string             `json:"id"`
	FullName       string             `json:"fullName"`
	Email          *string            `json:"email"`
	IdentityNumber string             `json:"identityNumber"`
	IdentityType   string             `json:"identityType"`
	PhoneNumber    *string            `json:"phoneNumber"`
	IsVerified     bool               `json:"isVerified"`
	AvatarURL      *string            `json:"avatarUrl"`
	Roles          []RoleResponse     `json:"roles"`
	Student        *StudentResponse   `json:"student,omitempty"`
	Lecturer       *LecturerResponse  `json:"lecturer,omitempty"`
	CreatedAt      string             `json:"createdAt"`
	UpdatedAt      string             `json:"updatedAt"`
}

type RoleResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type StudentResponse struct {
	ID             string `json:"id"`
	EnrollmentYear *int   `json:"enrollmentYear"`
	SKSCompleted   int    `json:"sksCompleted"`
	Status         string `json:"status"`
}

type LecturerResponse struct {
	ID             string  `json:"id"`
	ScienceGroupID *string `json:"scienceGroupId"`
}
