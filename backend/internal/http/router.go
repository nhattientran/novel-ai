package http

import (
	"novel-ai/internal/http/handlers"
	"novel-ai/internal/repo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router holds all HTTP routes
type Router struct {
	engine *gin.Engine
}

// NewRouter creates and configures the Gin router
func NewRouter(neo4j *repo.Driver) *Router {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Health handlers
	healthHandler := handlers.NewHealthHandler(neo4j)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Health)
		api.GET("/ready", healthHandler.Ready)
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
