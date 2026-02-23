package services

import (
	"context"
	"fmt"

	"novel-ai/internal/domain"
	"novel-ai/internal/repo"
)

// GraphService handles graph business logic
type GraphService struct {
	graphRepo *repo.GraphRepo
}

// NewGraphService creates a new graph service
func NewGraphService(graphRepo *repo.GraphRepo) *GraphService {
	return &GraphService{graphRepo: graphRepo}
}

// LoadGraph loads the complete story graph for Vue Flow
func (s *GraphService) LoadGraph(ctx context.Context, userID, storyID string) (*domain.GraphResponse, error) {
	graph, err := s.graphRepo.LoadGraph(ctx, userID, storyID)
	if err != nil {
		return nil, err
	}
	if graph == nil {
		return nil, fmt.Errorf("story not found")
	}
	return graph, nil
}
