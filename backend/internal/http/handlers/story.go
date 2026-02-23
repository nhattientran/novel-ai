package handlers

import (
	"net/http"

	"novel-ai/internal/auth"
	"novel-ai/internal/domain"
	"novel-ai/internal/repo"
	"novel-ai/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StoryHandler handles story CRUD endpoints
type StoryHandler struct {
	storyRepo *repo.StoryRepo
	storage   *storage.LocalStorage
}

// NewStoryHandler creates a new story handler
func NewStoryHandler(storyRepo *repo.StoryRepo, storage *storage.LocalStorage) *StoryHandler {
	return &StoryHandler{
		storyRepo: storyRepo,
		storage:   storage,
	}
}

// CreateStory handles story creation
func (h *StoryHandler) CreateStory(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	story := &domain.Story{
		ID:         uuid.New().String(),
		Title:      req.Title,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
	}

	createdStory, err := h.storyRepo.Create(c.Request.Context(), userID, story)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create story"})
		return
	}

	c.JSON(http.StatusCreated, createdStory)
}

// ListStories returns all stories for the current user
func (h *StoryHandler) ListStories(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stories, err := h.storyRepo.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list stories"})
		return
	}

	c.JSON(http.StatusOK, stories)
}

// GetStory returns a single story by ID
func (h *StoryHandler) GetStory(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	storyID := c.Param("storyId")
	if storyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID is required"})
		return
	}

	story, err := h.storyRepo.GetByID(c.Request.Context(), userID, storyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get story"})
		return
	}

	if story == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "story not found"})
		return
	}

	c.JSON(http.StatusOK, story)
}

// UpdateStory handles partial story updates
func (h *StoryHandler) UpdateStory(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	storyID := c.Param("storyId")
	if storyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID is required"})
		return
	}

	var req domain.UpdateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if at least one field is provided
	if req.Title == "" && req.Summary == "" && req.CoverImage == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one field is required for update"})
		return
	}

	updatedStory, err := h.storyRepo.Update(c.Request.Context(), userID, storyID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update story"})
		return
	}

	if updatedStory == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "story not found"})
		return
	}

	c.JSON(http.StatusOK, updatedStory)
}

// DeleteStory handles story deletion
func (h *StoryHandler) DeleteStory(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	storyID := c.Param("storyId")
	if storyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID is required"})
		return
	}

	// Get story first to check if it exists and get cover image
	story, err := h.storyRepo.GetByID(c.Request.Context(), userID, storyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get story"})
		return
	}

	if story == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "story not found"})
		return
	}

	// Delete the story (and its scenes)
	if err := h.storyRepo.Delete(c.Request.Context(), userID, storyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete story"})
		return
	}

	// Delete cover image if exists
	if story.CoverImage != "" {
		_ = h.storage.DeleteImage(story.CoverImage)
	}

	c.JSON(http.StatusOK, gin.H{"message": "story deleted successfully"})
}

// UploadImage handles image uploads for story covers
func (h *StoryHandler) UploadImage(c *gin.Context) {
	_, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	url, err := h.storage.UploadImage(file, header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
