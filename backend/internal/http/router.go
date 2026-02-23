package http

import (
	"novel-ai/internal/http/handlers"
	"novel-ai/internal/http/middleware"
	"novel-ai/internal/repo"
	"novel-ai/internal/services"
	"novel-ai/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router holds all HTTP routes
type Router struct {
	engine *gin.Engine
}

// RouterConfig holds configuration for router setup
type RouterConfig struct {
	Neo4j       *repo.Driver
	UserRepo    *repo.UserRepo
	StoryRepo   *repo.StoryRepo
	SceneRepo   *repo.SceneRepo
	ChoiceRepo  *repo.ChoiceRepo
	GraphRepo   *repo.GraphRepo
	Storage     *storage.LocalStorage
	JWTSecret   string
}

// NewRouter creates and configures the Gin router
func NewRouter(cfg *RouterConfig) *Router {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Ensure uploads directory exists
	cfg.Storage.EnsureUploadsDir()

	// Static file serving for uploads
	r.Static("/uploads", "./uploads")

	// Health handlers
	healthHandler := handlers.NewHealthHandler(cfg.Neo4j)

	// Auth handlers
	authHandler := handlers.NewAuthHandler(cfg.UserRepo, cfg.JWTSecret)

	// Story handlers
	storyHandler := handlers.NewStoryHandler(cfg.StoryRepo, cfg.Storage)

	// Scene service and handlers
	sceneService := services.NewSceneService(cfg.SceneRepo)
	sceneHandler := handlers.NewSceneHandler(sceneService)

	// Choice service and handlers
	choiceService := services.NewChoiceService(cfg.ChoiceRepo)
	choiceHandler := handlers.NewChoiceHandler(choiceService)

	// Graph service and handlers
	graphService := services.NewGraphService(cfg.GraphRepo)
	graphHandler := handlers.NewGraphHandler(graphService)

	// Auth middleware
	authMiddleware := middleware.Auth(cfg.JWTSecret)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Health)
		api.GET("/ready", healthHandler.Ready)

		// Auth routes (public)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/logout", authHandler.Logout)

		// Protected routes
		api.GET("/me", authMiddleware, authHandler.Me)

		// Creator routes
		creator := api.Group("/creator")
		creator.Use(authMiddleware, middleware.RequireRole("creator"))
		{
			// Story CRUD
			creator.POST("/stories", storyHandler.CreateStory)
			creator.GET("/stories", storyHandler.ListStories)
			creator.GET("/stories/:storyId", storyHandler.GetStory)
			creator.PATCH("/stories/:storyId", storyHandler.UpdateStory)
			creator.DELETE("/stories/:storyId", storyHandler.DeleteStory)

			// Scene CRUD
			creator.POST("/stories/:storyId/scenes", sceneHandler.CreateScene)
			creator.GET("/stories/:storyId/scenes/:sceneId", sceneHandler.GetScene)
			creator.PATCH("/stories/:storyId/scenes/:sceneId", sceneHandler.UpdateScene)
			creator.DELETE("/stories/:storyId/scenes/:sceneId", sceneHandler.DeleteScene)

			// Start scene
			creator.PUT("/stories/:storyId/start/:sceneId", sceneHandler.SetStartScene)

			// Choice management
			creator.POST("/stories/:storyId/choices", choiceHandler.CreateChoice)
			creator.PATCH("/stories/:storyId/choices", choiceHandler.UpdateChoice)
			creator.DELETE("/stories/:storyId/choices", choiceHandler.DeleteChoice)

			// Graph loading
			creator.GET("/stories/:storyId/graph", graphHandler.LoadGraph)
		}

		// Upload routes (requires auth)
		api.POST("/uploads/images", authMiddleware, storyHandler.UploadImage)
	}

	return &Router{engine: r}
}

// Engine returns the gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// Run starts the HTTP server
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
