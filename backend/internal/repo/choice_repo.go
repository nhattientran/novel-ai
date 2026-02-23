package repo

import (
	"context"
	"fmt"

	"novel-ai/internal/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ChoiceRepo handles choice (edge) data access
type ChoiceRepo struct {
	driver *Driver
}

// NewChoiceRepo creates a new choice repository
func NewChoiceRepo(driver *Driver) *ChoiceRepo {
	return &ChoiceRepo{driver: driver}
}

// Upsert creates or updates a choice between two scenes
func (r *ChoiceRepo) Upsert(ctx context.Context, userID, storyID string, choice *domain.Choice) (*domain.Choice, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})
		MATCH (st)-[:HAS_SCENE]->(a:Scene {id: $from_scene_id})
		MATCH (st)-[:HAS_SCENE]->(b:Scene {id: $to_scene_id})
		MERGE (a)-[r:LEADS_TO]->(b)
		SET r.choice_text = $choice_text
		RETURN a.id AS from, b.id AS to, r.choice_text AS choice_text
	`

	params := map[string]any{
		"user_id":        userID,
		"story_id":       storyID,
		"from_scene_id":  choice.FromSceneID,
		"to_scene_id":    choice.ToSceneID,
		"choice_text":    choice.ChoiceText,
	}

	records, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert choice: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("scenes not found or not in story")
	}

	return r.recordToChoice(records[0])
}

// Delete removes a choice between two scenes
func (r *ChoiceRepo) Delete(ctx context.Context, userID, storyID, fromSceneID, toSceneID string) error {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})
		MATCH (st)-[:HAS_SCENE]->(a:Scene {id: $from_scene_id})-[r:LEADS_TO]->(b:Scene {id: $to_scene_id})<-[:HAS_SCENE]-(st)
		DELETE r
	`

	params := map[string]any{
		"user_id":        userID,
		"story_id":       storyID,
		"from_scene_id":  fromSceneID,
		"to_scene_id":    toSceneID,
	}

	_, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete choice: %w", err)
	}

	return nil
}

// ListByStory returns all choices for a story
func (r *ChoiceRepo) ListByStory(ctx context.Context, userID, storyID string) ([]*domain.Choice, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})
		MATCH (st)-[:HAS_SCENE]->(a:Scene)-[r:LEADS_TO]->(b:Scene)<-[:HAS_SCENE]-(st)
		RETURN a.id AS from, b.id AS to, r.choice_text AS choice_text
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list choices: %w", err)
	}

	choices := make([]*domain.Choice, 0, len(records))
	for _, record := range records {
		choice, err := r.recordToChoice(record)
		if err != nil {
			return nil, err
		}
		choices = append(choices, choice)
	}

	return choices, nil
}

// recordToChoice converts a Neo4j record to Choice domain model
func (r *ChoiceRepo) recordToChoice(record *neo4j.Record) (*domain.Choice, error) {
	from, ok := record.Get("from")
	if !ok {
		return nil, fmt.Errorf("from not found in record")
	}

	to, ok := record.Get("to")
	if !ok {
		return nil, fmt.Errorf("to not found in record")
	}

	choiceText, ok := record.Get("choice_text")
	if !ok {
		return nil, fmt.Errorf("choice_text not found in record")
	}

	return &domain.Choice{
		FromSceneID: from.(string),
		ToSceneID:   to.(string),
		ChoiceText:  choiceText.(string),
	}, nil
}
