package handlers

import (
	"log"
	"net/http"

	"novel-ai/internal/auth"
	"novel-ai/internal/domain"
	"novel-ai/internal/services"

	"github.com/gin-gonic/gin"
)

// ChoiceHandler handles choice (edge) endpoints
type ChoiceHandler struct {
	choiceService *services.ChoiceService
}

// NewChoiceHandler creates a new choice handler
func NewChoiceHandler(choiceService *services.ChoiceService) *ChoiceHandler {
	return &ChoiceHandler{choiceService: choiceService}
}

// CreateChoice creates or updates a choice between two scenes
func (h *ChoiceHandler) CreateChoice(c *gin.Context) {
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

	var req domain.CreateChoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdChoice, err := h.choiceService.CreateOrUpdate(c.Request.Context(), userID, storyID, &req)
	if err != nil {
		log.Printf("Failed to create choice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create choice"})
		return
	}

	c.JSON(http.StatusCreated, createdChoice)
}

// UpdateChoice updates a choice's text
func (h *ChoiceHandler) UpdateChoice(c *gin.Context) {
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

	var req domain.UpdateChoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedChoice, err := h.choiceService.Update(c.Request.Context(), userID, storyID, &req)
	if err != nil {
		log.Printf("Failed to update choice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update choice"})
		return
	}

	c.JSON(http.StatusOK, updatedChoice)
}

// DeleteChoice deletes a choice between two scenes
func (h *ChoiceHandler) DeleteChoice(c *gin.Context) {
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

	var req domain.DeleteChoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.choiceService.Delete(c.Request.Context(), userID, storyID, &req); err != nil {
		log.Printf("Failed to delete choice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete choice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "choice deleted successfully"})
}
