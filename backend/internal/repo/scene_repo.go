package repo

import (
	"context"
	"fmt"

	"novel-ai/internal/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// SceneRepo handles scene data access
type SceneRepo struct {
	driver *Driver
}

// NewSceneRepo creates a new scene repository
func NewSceneRepo(driver *Driver) *SceneRepo {
	return &SceneRepo{driver: driver}
}

// Create creates a new scene for a story
func (r *SceneRepo) Create(ctx context.Context, userID, storyID string, scene *domain.Scene) (*domain.Scene, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})
		CREATE (sc:Scene {
			id: $scene_id,
			content: $content,
			image_url: $image_url,
			is_end: coalesce($is_end, false),
			pos_x: coalesce($pos_x, 0.0),
			pos_y: coalesce($pos_y, 0.0)
		})
		CREATE (st)-[:HAS_SCENE]->(sc)
		RETURN sc { .id, .content, .image_url, .is_end, .pos_x, .pos_y } AS scene
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
		"scene_id":  scene.ID,
		"content":   scene.Content,
		"image_url": scene.ImageURL,
		"is_end":    scene.IsEnd,
		"pos_x":     scene.PosX,
		"pos_y":     scene.PosY,
	}

	records, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create scene: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no scene created")
	}

	return r.recordToScene(records[0])
}

// GetByID returns a scene by ID if owned by the user
func (r *SceneRepo) GetByID(ctx context.Context, userID, storyID, sceneID string) (*domain.Scene, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})-[:HAS_SCENE]->(sc:Scene {id: $scene_id})
		RETURN sc { .id, .content, .image_url, .is_end, .pos_x, .pos_y } AS scene
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
		"scene_id":  sceneID,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get scene: %w", err)
	}

	if len(records) == 0 {
		return nil, nil // scene not found or not owned by user
	}

	return r.recordToScene(records[0])
}

// Update updates a scene (partial update using coalesce)
func (r *SceneRepo) Update(ctx context.Context, userID, storyID, sceneID string, req *domain.UpdateSceneRequest) (*domain.Scene, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})-[:HAS_SCENE]->(sc:Scene {id: $scene_id})
		SET sc.content = coalesce($content, sc.content),
		    sc.image_url = coalesce($image_url, sc.image_url),
		    sc.is_end = coalesce($is_end, sc.is_end),
		    sc.pos_x = coalesce($pos_x, sc.pos_x),
		    sc.pos_y = coalesce($pos_y, sc.pos_y)
		RETURN sc { .id, .content, .image_url, .is_end, .pos_x, .pos_y } AS scene
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
		"scene_id":  sceneID,
		"content":   req.Content,
		"image_url": req.ImageURL,
		"is_end":    req.IsEnd,
		"pos_x":     req.PosX,
		"pos_y":     req.PosY,
	}

	records, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update scene: %w", err)
	}

	if len(records) == 0 {
		return nil, nil // scene not found or not owned by user
	}

	return r.recordToScene(records[0])
}

// Delete deletes a scene (and all its relationships)
func (r *SceneRepo) Delete(ctx context.Context, userID, storyID, sceneID string) error {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})-[:HAS_SCENE]->(sc:Scene {id: $scene_id})
		DETACH DELETE sc
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
		"scene_id":  sceneID,
	}

	_, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete scene: %w", err)
	}

	return nil
}

// ListByStory returns all scenes for a story
func (r *SceneRepo) ListByStory(ctx context.Context, userID, storyID string) ([]*domain.Scene, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})-[:HAS_SCENE]->(sc:Scene)
		RETURN sc { .id, .content, .image_url, .is_end, .pos_x, .pos_y } AS scene
		ORDER BY sc.id
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list scenes: %w", err)
	}

	scenes := make([]*domain.Scene, 0, len(records))
	for _, record := range records {
		scene, err := r.recordToScene(record)
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, scene)
	}

	return scenes, nil
}

// SetStartScene sets a scene as the start scene for a story
// This removes any existing STARTS_AT relationship first
func (r *SceneRepo) SetStartScene(ctx context.Context, userID, storyID, sceneID string) error {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})
		OPTIONAL MATCH (st)-[old:STARTS_AT]->(:Scene)
		DELETE old
		WITH st
		MATCH (st)-[:HAS_SCENE]->(sc:Scene {id: $scene_id})
		MERGE (st)-[:STARTS_AT]->(sc)
		RETURN sc.id AS start_scene_id
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
		"scene_id":  sceneID,
	}

	records, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to set start scene: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("scene not found or not owned by user")
	}

	return nil
}

// recordToScene converts a Neo4j record to Scene domain model
func (r *SceneRepo) recordToScene(record *neo4j.Record) (*domain.Scene, error) {
	sceneMap, ok := record.Get("scene")
	if !ok {
		return nil, fmt.Errorf("scene not found in record")
	}

	sceneData, ok := sceneMap.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid scene data type")
	}

	scene := &domain.Scene{
		ID:      getString(sceneData, "id"),
		Content: getString(sceneData, "content"),
		IsEnd:   getBool(sceneData, "is_end"),
		PosX:    getFloat64(sceneData, "pos_x"),
		PosY:    getFloat64(sceneData, "pos_y"),
	}

	// Handle nullable image_url
	if imageURL, ok := sceneData["image_url"]; ok && imageURL != nil {
		if s, ok := imageURL.(string); ok && s != "" {
			scene.ImageURL = &s
		}
	}

	return scene, nil
}
