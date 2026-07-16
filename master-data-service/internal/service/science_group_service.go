package service

import (
	"context"

	"github.com/google/uuid"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/model"
	"neocentral-go/master-data-service/internal/repository"
	"neocentral-go/pkg/apperror"
)

type ScienceGroupService struct {
	repo repository.ScienceGroupRepository
}

func NewScienceGroupService(repo repository.ScienceGroupRepository) *ScienceGroupService {
	return &ScienceGroupService{repo: repo}
}

func (s *ScienceGroupService) GetAll(ctx context.Context) ([]dto.ScienceGroupResponse, error) {
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperror.InternalWrap("failed to fetch science groups", err)
	}
	result := make([]dto.ScienceGroupResponse, len(list))
	for i, sg := range list {
		result[i] = toScienceGroupResponse(sg)
	}
	return result, nil
}

func (s *ScienceGroupService) GetByID(ctx context.Context, id string) (*dto.ScienceGroupResponse, error) {
	sg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if sg == nil {
		return nil, apperror.NotFound("Science group not found")
	}
	resp := toScienceGroupResponse(*sg)
	return &resp, nil
}

func (s *ScienceGroupService) Create(ctx context.Context, req dto.CreateScienceGroupRequest) (*dto.ScienceGroupResponse, error) {
	sg := model.ScienceGroup{
		ID:   uuid.NewString(),
		Name: req.Name,
	}
	if err := s.repo.Create(ctx, &sg); err != nil {
		return nil, apperror.InternalWrap("failed to create science group", err)
	}
	resp := toScienceGroupResponse(sg)
	return &resp, nil
}

func (s *ScienceGroupService) Update(ctx context.Context, id string, req dto.UpdateScienceGroupRequest) (*dto.ScienceGroupResponse, error) {
	sg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if sg == nil {
		return nil, apperror.NotFound("Science group not found")
	}
	if req.Name != nil {
		sg.Name = *req.Name
	}
	if err := s.repo.Update(ctx, sg); err != nil {
		return nil, apperror.InternalWrap("failed to update science group", err)
	}
	resp := toScienceGroupResponse(*sg)
	return &resp, nil
}

func (s *ScienceGroupService) Delete(ctx context.Context, id string) error {
	sg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperror.InternalWrap("database error", err)
	}
	if sg == nil {
		return apperror.NotFound("Science group not found")
	}
	return s.repo.Delete(ctx, id)
}

func toScienceGroupResponse(sg model.ScienceGroup) dto.ScienceGroupResponse {
	return dto.ScienceGroupResponse{
		ID:        sg.ID,
		Name:      sg.Name,
		CreatedAt: sg.CreatedAt,
		UpdatedAt: sg.UpdatedAt,
	}
}
