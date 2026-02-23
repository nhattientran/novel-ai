package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"novel-ai/internal/repo"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	neo4j *repo.Driver
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(neo4j *repo.Driver) *HealthHandler {
	return &HealthHandler{neo4j: neo4j}
}

// Health returns basic health status
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Ready checks if all dependencies are ready
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.neo4j.VerifyConnectivity(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"error":   "neo4j connection failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
