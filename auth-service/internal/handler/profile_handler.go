package handler

import (
	"net/http"

	"neocentral-go/auth-service/internal/service"
	"neocentral-go/pkg/response"
)

// ProfileHandler exposes HTTP endpoints for user profiles.
type ProfileHandler struct {
	profileSvc *service.ProfileService
}

func NewProfileHandler(svc *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileSvc: svc}
}

// GetMe retrieves the authenticated user's profile.
// @Summary      Get My Profile
// @Description  Retrieves the profile information of the currently authenticated user.
// @Tags         profile
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.APIResponse{data=dto.UserProfileResponse}
// @Failure      401  {object}  response.APIResponse
// @Router       /profile/me [get]
func (h *ProfileHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	profile, err := h.profileSvc.GetProfile(r.Context(), userID)
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}
