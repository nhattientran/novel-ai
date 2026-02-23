package services

import (
	"context"

	"novel-ai/internal/domain"
	"novel-ai/internal/repo"
)

// ChoiceService handles choice business logic
type ChoiceService struct {
	choiceRepo *repo.ChoiceRepo
}

// NewChoiceService creates a new choice service
func NewChoiceService(choiceRepo *repo.ChoiceRepo) *ChoiceService {
	return &ChoiceService{choiceRepo: choiceRepo}
}

// CreateOrUpdate creates or updates a choice between two scenes
func (s *ChoiceService) CreateOrUpdate(ctx context.Context, userID, storyID string, req *domain.CreateChoiceRequest) (*domain.Choice, error) {
	choice := &domain.Choice{
		FromSceneID: req.FromSceneID,
		ToSceneID:   req.ToSceneID,
		ChoiceText:  req.ChoiceText,
	}

	return s.choiceRepo.Upsert(ctx, userID, storyID, choice)
}

// Update updates a choice's text
func (s *ChoiceService) Update(ctx context.Context, userID, storyID string, req *domain.UpdateChoiceRequest) (*domain.Choice, error) {
	choice := &domain.Choice{
		FromSceneID: req.FromSceneID,
		ToSceneID:   req.ToSceneID,
		ChoiceText:  req.ChoiceText,
	}

	return s.choiceRepo.Upsert(ctx, userID, storyID, choice)
}

// Delete deletes a choice
func (s *ChoiceService) Delete(ctx context.Context, userID, storyID string, req *domain.DeleteChoiceRequest) error {
	return s.choiceRepo.Delete(ctx, userID, storyID, req.FromSceneID, req.ToSceneID)
}

// List returns all choices for a story
func (s *ChoiceService) List(ctx context.Context, userID, storyID string) ([]*domain.Choice, error) {
	return s.choiceRepo.ListByStory(ctx, userID, storyID)
}
