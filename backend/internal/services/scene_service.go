package services

import (
	"context"
	"fmt"

	"novel-ai/internal/domain"
	"novel-ai/internal/repo"

	"github.com/google/uuid"
)

// SceneService handles scene business logic
type SceneService struct {
	sceneRepo *repo.SceneRepo
}

// NewSceneService creates a new scene service
func NewSceneService(sceneRepo *repo.SceneRepo) *SceneService {
	return &SceneService{sceneRepo: sceneRepo}
}

// Create creates a new scene for a story
func (s *SceneService) Create(ctx context.Context, userID, storyID string, req *domain.CreateSceneRequest) (*domain.Scene, error) {
	scene := &domain.Scene{
		ID:       uuid.New().String(),
		Content:  req.Content,
		ImageURL: req.ImageURL,
		IsEnd:    req.IsEnd,
		PosX:     req.PosX,
		PosY:     req.PosY,
	}

	return s.sceneRepo.Create(ctx, userID, storyID, scene)
}

// Get returns a scene by ID
func (s *SceneService) Get(ctx context.Context, userID, storyID, sceneID string) (*domain.Scene, error) {
	scene, err := s.sceneRepo.GetByID(ctx, userID, storyID, sceneID)
	if err != nil {
		return nil, err
	}
	if scene == nil {
		return nil, fmt.Errorf("scene not found")
	}
	return scene, nil
}

// Update updates a scene
func (s *SceneService) Update(ctx context.Context, userID, storyID, sceneID string, req *domain.UpdateSceneRequest) (*domain.Scene, error) {
	scene, err := s.sceneRepo.Update(ctx, userID, storyID, sceneID, req)
	if err != nil {
		return nil, err
	}
	if scene == nil {
		return nil, fmt.Errorf("scene not found")
	}
	return scene, nil
}

// Delete deletes a scene
func (s *SceneService) Delete(ctx context.Context, userID, storyID, sceneID string) error {
	return s.sceneRepo.Delete(ctx, userID, storyID, sceneID)
}

// List returns all scenes for a story
func (s *SceneService) List(ctx context.Context, userID, storyID string) ([]*domain.Scene, error) {
	return s.sceneRepo.ListByStory(ctx, userID, storyID)
}

// SetStartScene sets a scene as the start scene for a story
func (s *SceneService) SetStartScene(ctx context.Context, userID, storyID, sceneID string) error {
	return s.sceneRepo.SetStartScene(ctx, userID, storyID, sceneID)
}
