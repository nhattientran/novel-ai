package handlers

import (
	"log"
	"net/http"

	"novel-ai/internal/auth"
	"novel-ai/internal/domain"
	"novel-ai/internal/services"

	"github.com/gin-gonic/gin"
)

// SceneHandler handles scene CRUD endpoints
type SceneHandler struct {
	sceneService *services.SceneService
}

// NewSceneHandler creates a new scene handler
func NewSceneHandler(sceneService *services.SceneService) *SceneHandler {
	return &SceneHandler{sceneService: sceneService}
}

// CreateScene handles scene creation
func (h *SceneHandler) CreateScene(c *gin.Context) {
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

	var req domain.CreateSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdScene, err := h.sceneService.Create(c.Request.Context(), userID, storyID, &req)
	if err != nil {
		log.Printf("Failed to create scene: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create scene"})
		return
	}

	c.JSON(http.StatusCreated, createdScene)
}

// GetScene returns a single scene by ID
func (h *SceneHandler) GetScene(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	storyID := c.Param("storyId")
	sceneID := c.Param("sceneId")
	if storyID == "" || sceneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID and scene ID are required"})
		return
	}

	scene, err := h.sceneService.Get(c.Request.Context(), userID, storyID, sceneID)
	if err != nil {
		log.Printf("Failed to get scene: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get scene"})
		return
	}

	if scene == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scene not found"})
		return
	}

	c.JSON(http.StatusOK, scene)
}

// UpdateScene handles partial scene updates
func (h *SceneHandler) UpdateScene(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	storyID := c.Param("storyId")
	sceneID := c.Param("sceneId")
	if storyID == "" || sceneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID and scene ID are required"})
		return
	}

	var req domain.UpdateSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedScene, err := h.sceneService.Update(c.Request.Context(), userID, storyID, sceneID, &req)
	if err != nil {
		log.Printf("Failed to update scene: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update scene"})
		return
	}

	c.JSON(http.StatusOK, updatedScene)
}

// DeleteScene handles scene deletion
func (h *SceneHandler) DeleteScene(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	storyID := c.Param("storyId")
	sceneID := c.Param("sceneId")
	if storyID == "" || sceneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID and scene ID are required"})
		return
	}

	if err := h.sceneService.Delete(c.Request.Context(), userID, storyID, sceneID); err != nil {
		log.Printf("Failed to delete scene: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete scene"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "scene deleted successfully"})
}

// SetStartScene sets a scene as the start scene for a story
func (h *SceneHandler) SetStartScene(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	storyID := c.Param("storyId")
	sceneID := c.Param("sceneId")
	if storyID == "" || sceneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "story ID and scene ID are required"})
		return
	}

	if err := h.sceneService.SetStartScene(c.Request.Context(), userID, storyID, sceneID); err != nil {
		log.Printf("Failed to set start scene: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set start scene"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "start scene set successfully", "start_scene_id": sceneID})
}
