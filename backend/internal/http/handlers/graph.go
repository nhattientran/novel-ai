package handlers

import (
	"log"
	"net/http"

	"novel-ai/internal/auth"
	"novel-ai/internal/services"

	"github.com/gin-gonic/gin"
)

// GraphHandler handles graph loading endpoints
type GraphHandler struct {
	graphService *services.GraphService
}

// NewGraphHandler creates a new graph handler
func NewGraphHandler(graphService *services.GraphService) *GraphHandler {
	return &GraphHandler{graphService: graphService}
}

// LoadGraph loads the complete story graph for Vue Flow
func (h *GraphHandler) LoadGraph(c *gin.Context) {
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

	log.Printf("Loading graph for userID=%s, storyID=%s", userID, storyID)

	graph, err := h.graphService.LoadGraph(c.Request.Context(), userID, storyID)
	if err != nil {
		log.Printf("Failed to load graph: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if graph == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "story not found"})
		return
	}

	c.JSON(http.StatusOK, graph)
}
