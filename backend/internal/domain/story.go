package domain

import "time"

// Story represents a light novel story created by a user
type Story struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	CoverImage string    `json:"cover_image"`
	Status     string    `json:"status"` // "draft" or "published"
	CreatedAt  time.Time `json:"created_at"`
}

// CreateStoryRequest represents the request to create a new story
type CreateStoryRequest struct {
	Title      string `json:"title" binding:"required,max=200"`
	Summary    string `json:"summary" binding:"max=2000"`
	CoverImage string `json:"cover_image"`
}

// UpdateStoryRequest represents the request to update a story
type UpdateStoryRequest struct {
	Title      string `json:"title" binding:"omitempty,max=200"`
	Summary    string `json:"summary" binding:"omitempty,max=2000"`
	CoverImage string `json:"cover_image"`
}
