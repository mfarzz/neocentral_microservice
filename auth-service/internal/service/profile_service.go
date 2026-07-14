package service

import (
	"context"

	"neocentral-go/auth-service/internal/dto"
	"neocentral-go/auth-service/internal/repository"
	"neocentral-go/pkg/apperror"
)

// ProfileService handles user profile queries.
type ProfileService struct {
	userRepo repository.UserRepository
}

func NewProfileService(repo repository.UserRepository) *ProfileService {
	return &ProfileService{userRepo: repo}
}

// GetProfile returns the full profile for the given user ID.
func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}

	profile := buildProfileResponse(user)
	return &profile, nil
}
