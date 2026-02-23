package repo

import (
	"context"
	"fmt"
	"time"

	"novel-ai/internal/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// StoryRepo handles story data access
type StoryRepo struct {
	driver *Driver
}

// NewStoryRepo creates a new story repository
func NewStoryRepo(driver *Driver) *StoryRepo {
	return &StoryRepo{driver: driver}
}

// Create creates a new story for a user
func (r *StoryRepo) Create(ctx context.Context, userID string, story *domain.Story) (*domain.Story, error) {
	query := `
		MATCH (u:User {id: $user_id})
		CREATE (s:Story {
			id: $story_id,
			title: $title,
			summary: $summary,
			cover_image: $cover_image,
			status: "draft",
			created_at: datetime()
		})
		CREATE (u)-[:CREATED]->(s)
		RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story
	`

	params := map[string]any{
		"user_id":     userID,
		"story_id":    story.ID,
		"title":       story.Title,
		"summary":     story.Summary,
		"cover_image": story.CoverImage,
	}

	records, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create story: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no story created")
	}

	return r.recordToStory(records[0])
}

// ListByUser returns all stories created by a user
func (r *StoryRepo) ListByUser(ctx context.Context, userID string) ([]*domain.Story, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story)
		RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story
		ORDER BY s.created_at DESC
	`

	params := map[string]any{
		"user_id": userID,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list stories: %w", err)
	}

	stories := make([]*domain.Story, 0, len(records))
	for _, record := range records {
		story, err := r.recordToStory(record)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	return stories, nil
}

// GetByID returns a story by ID if owned by the user
func (r *StoryRepo) GetByID(ctx context.Context, userID, storyID string) (*domain.Story, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story {id: $story_id})
		RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story
	`

	params := map[string]any{
		"user_id":  userID,
		"story_id": storyID,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get story: %w", err)
	}

	if len(records) == 0 {
		return nil, nil // story not found or not owned by user
	}

	return r.recordToStory(records[0])
}

// Update updates a story (partial update using coalesce)
func (r *StoryRepo) Update(ctx context.Context, userID, storyID string, req *domain.UpdateStoryRequest) (*domain.Story, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story {id: $story_id})
		SET s.title = coalesce($title, s.title),
		    s.summary = coalesce($summary, s.summary),
		    s.cover_image = coalesce($cover_image, s.cover_image)
		RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story
	`

	params := map[string]any{
		"user_id":     userID,
		"story_id":    storyID,
		"title":       req.Title,
		"summary":     req.Summary,
		"cover_image": req.CoverImage,
	}

	records, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update story: %w", err)
	}

	if len(records) == 0 {
		return nil, nil // story not found or not owned by user
	}

	return r.recordToStory(records[0])
}

// Delete deletes a story and all its scenes (if owned by user)
func (r *StoryRepo) Delete(ctx context.Context, userID, storyID string) error {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story {id: $story_id})
		OPTIONAL MATCH (s)-[:HAS_SCENE]->(sc:Scene)
		DETACH DELETE sc, s
	`

	params := map[string]any{
		"user_id":  userID,
		"story_id": storyID,
	}

	_, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete story: %w", err)
	}

	return nil
}

// recordToStory converts a Neo4j record to Story domain model
func (r *StoryRepo) recordToStory(record *neo4j.Record) (*domain.Story, error) {
	storyMap, ok := record.Get("story")
	if !ok {
		return nil, fmt.Errorf("story not found in record")
	}

	storyData, ok := storyMap.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid story data type")
	}

	story := &domain.Story{
		ID:         getString(storyData, "id"),
		Title:      getString(storyData, "title"),
		Summary:    getString(storyData, "summary"),
		CoverImage: getString(storyData, "cover_image"),
		Status:     getString(storyData, "status"),
	}

	// Parse created_at time
	if createdAt, ok := storyData["created_at"]; ok {
		switch t := createdAt.(type) {
		case time.Time:
			story.CreatedAt = t
		case neo4j.Date:
			story.CreatedAt = t.Time()
		case neo4j.LocalDateTime:
			story.CreatedAt = t.Time()
		case neo4j.OffsetTime:
			story.CreatedAt = t.Time()
		}
	}

	return story, nil
}
