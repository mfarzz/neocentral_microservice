package service

import (
	"context"

	"github.com/google/uuid"
	"neocentral-go/master-data-service/internal/dto"
	"neocentral-go/master-data-service/internal/model"
	"neocentral-go/master-data-service/internal/repository"
	"neocentral-go/pkg/apperror"
)

type ThesisTopicService struct {
	topicRepo  repository.ThesisTopicRepository
	statusRepo repository.ThesisStatusRepository
}

func NewThesisTopicService(topicRepo repository.ThesisTopicRepository, statusRepo repository.ThesisStatusRepository) *ThesisTopicService {
	return &ThesisTopicService{topicRepo: topicRepo, statusRepo: statusRepo}
}

// ── Thesis Topics ────────────────────────────

func (s *ThesisTopicService) GetAllTopics(ctx context.Context) ([]dto.ThesisTopicResponse, error) {
	list, err := s.topicRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.InternalWrap("failed to fetch thesis topics", err)
	}
	result := make([]dto.ThesisTopicResponse, len(list))
	for i, t := range list {
		result[i] = toThesisTopicResponse(t)
	}
	return result, nil
}

func (s *ThesisTopicService) GetTopicByID(ctx context.Context, id string) (*dto.ThesisTopicResponse, error) {
	topic, err := s.topicRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if topic == nil {
		return nil, apperror.NotFound("Thesis topic not found")
	}
	resp := toThesisTopicResponse(*topic)
	return &resp, nil
}

func (s *ThesisTopicService) CreateTopic(ctx context.Context, req dto.CreateThesisTopicRequest) (*dto.ThesisTopicResponse, error) {
	topic := model.ThesisTopic{
		ID:   uuid.NewString(),
		Name: req.Name,
	}
	if err := s.topicRepo.Create(ctx, &topic); err != nil {
		return nil, apperror.InternalWrap("failed to create thesis topic", err)
	}
	resp := toThesisTopicResponse(topic)
	return &resp, nil
}

func (s *ThesisTopicService) UpdateTopic(ctx context.Context, id string, req dto.UpdateThesisTopicRequest) (*dto.ThesisTopicResponse, error) {
	topic, err := s.topicRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap("database error", err)
	}
	if topic == nil {
		return nil, apperror.NotFound("Thesis topic not found")
	}
	if req.Name != nil {
		topic.Name = *req.Name
	}
	if err := s.topicRepo.Update(ctx, topic); err != nil {
		return nil, apperror.InternalWrap("failed to update thesis topic", err)
	}
	resp := toThesisTopicResponse(*topic)
	return &resp, nil
}

func (s *ThesisTopicService) DeleteTopic(ctx context.Context, id string) error {
	topic, err := s.topicRepo.FindByID(ctx, id)
	if err != nil {
		return apperror.InternalWrap("database error", err)
	}
	if topic == nil {
		return apperror.NotFound("Thesis topic not found")
	}
	return s.topicRepo.Delete(ctx, id)
}

// ── Thesis Statuses ──────────────────────────

func (s *ThesisTopicService) GetAllStatuses(ctx context.Context) ([]dto.ThesisStatusResponse, error) {
	list, err := s.statusRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.InternalWrap("failed to fetch thesis statuses", err)
	}
	result := make([]dto.ThesisStatusResponse, len(list))
	for i, st := range list {
		result[i] = dto.ThesisStatusResponse{ID: st.ID, Name: st.Name}
	}
	return result, nil
}

func toThesisTopicResponse(t model.ThesisTopic) dto.ThesisTopicResponse {
	return dto.ThesisTopicResponse{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
