package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"novel-ai/internal/config"
	httpServer "novel-ai/internal/http"
	"novel-ai/internal/repo"
	"novel-ai/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Neo4j driver
	neo4jDriver, err := repo.NewDriver(cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword)
	if err != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", err)
	}

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := neo4jDriver.VerifyConnectivity(ctx); err != nil {
		log.Printf("Warning: Neo4j connectivity check failed: %v", err)
	} else {
		log.Println("Connected to Neo4j successfully")
	}

	// Create repositories
	userRepo := repo.NewUserRepo(neo4jDriver)
	storyRepo := repo.NewStoryRepo(neo4jDriver)
	sceneRepo := repo.NewSceneRepo(neo4jDriver)
	choiceRepo := repo.NewChoiceRepo(neo4jDriver)
	graphRepo := repo.NewGraphRepo(neo4jDriver)

	// Create storage
	localStorage := storage.NewLocalStorage("./uploads", "")

	// Create router config
	routerCfg := &httpServer.RouterConfig{
		Neo4j:      neo4jDriver,
		UserRepo:   userRepo,
		StoryRepo:  storyRepo,
		SceneRepo:  sceneRepo,
		ChoiceRepo: choiceRepo,
		GraphRepo:  graphRepo,
		Storage:    localStorage,
		JWTSecret:  cfg.JWTSecret,
	}

	// Create router
	router := httpServer.NewRouter(routerCfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router.Engine(),
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close Neo4j driver
	if err := neo4jDriver.Close(shutdownCtx); err != nil {
		log.Printf("Error closing Neo4j driver: %v", err)
	}

	log.Println("Server exited")
}
