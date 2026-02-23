package repo

import (
	"context"
	"fmt"

	"novel-ai/internal/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// GraphRepo handles loading the full story graph for Vue Flow
type GraphRepo struct {
	driver *Driver
}

// NewGraphRepo creates a new graph repository
func NewGraphRepo(driver *Driver) *GraphRepo {
	return &GraphRepo{driver: driver}
}

// LoadGraph loads the complete story graph including scenes and choices
func (r *GraphRepo) LoadGraph(ctx context.Context, userID, storyID string) (*domain.GraphResponse, error) {
	query := `
		MATCH (u:User {id: $user_id})-[:CREATED]->(st:Story {id: $story_id})
		OPTIONAL MATCH (st)-[:STARTS_AT]->(start:Scene)
		OPTIONAL MATCH (st)-[:HAS_SCENE]->(sc:Scene)
		OPTIONAL MATCH (sc)-[r:LEADS_TO]->(to:Scene)<-[:HAS_SCENE]-(st)
		RETURN
			st { .id, .title, .status } AS story,
			collect(DISTINCT sc {
				.id, .pos_x, .pos_y, .is_end, .content, .image_url,
				is_start: (start IS NOT NULL AND sc.id = start.id)
			}) AS scenes,
			collect(DISTINCT CASE WHEN to IS NOT NULL THEN {
				from: sc.id, to: to.id, choice_text: r.choice_text
			} END) AS choices
	`

	params := map[string]any{
		"user_id":   userID,
		"story_id":  storyID,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to load graph: %w", err)
	}

	if len(records) == 0 {
		// Try to determine if story exists but user doesn't own it
		checkQuery := `
			MATCH (st:Story {id: $story_id})
			RETURN st.id as story_id
		`
		checkRecords, checkErr := r.driver.ExecuteRead(ctx, checkQuery, map[string]any{"story_id": storyID})
		if checkErr != nil {
			return nil, fmt.Errorf("story not found: %w", checkErr)
		}
		if len(checkRecords) == 0 {
			return nil, fmt.Errorf("story not found: %s", storyID)
		}
		return nil, fmt.Errorf("story not owned by user: %s", storyID)
	}

	return r.recordToGraphResponse(records[0])
}

// recordToGraphResponse converts a Neo4j record to GraphResponse
func (r *GraphRepo) recordToGraphResponse(record *neo4j.Record) (*domain.GraphResponse, error) {
	// Parse story
	storyMap, ok := record.Get("story")
	if !ok {
		return nil, fmt.Errorf("story not found in record")
	}

	storyData, ok := storyMap.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid story data type")
	}

	story := domain.StorySummary{
		ID:     getString(storyData, "id"),
		Title:  getString(storyData, "title"),
		Status: getString(storyData, "status"),
	}

	// Parse scenes
	scenesList, ok := record.Get("scenes")
	if !ok {
		return nil, fmt.Errorf("scenes not found in record")
	}

	scenesData, ok := scenesList.([]any)
	if !ok {
		// Handle case where scenes is null (no scenes yet)
		scenesData = []any{}
	}

	nodes := make([]domain.GraphNode, 0, len(scenesData))
	for _, sceneItem := range scenesData {
		sceneMap, ok := sceneItem.(map[string]any)
		if !ok {
			continue
		}

		node := domain.GraphNode{
			ID:       getString(sceneMap, "id"),
			Type:     "scene",
			Position: domain.Position{
				X: getFloat64(sceneMap, "pos_x"),
				Y: getFloat64(sceneMap, "pos_y"),
			},
			Data: domain.NodeData{
				Content: getString(sceneMap, "content"),
				IsStart: getBool(sceneMap, "is_start"),
				IsEnd:   getBool(sceneMap, "is_end"),
			},
		}

		// Handle nullable image_url
		if imageURL, ok := sceneMap["image_url"]; ok && imageURL != nil {
			if s, ok := imageURL.(string); ok && s != "" {
				node.Data.ImageURL = &s
			}
		}

		nodes = append(nodes, node)
	}

	// Parse choices
	choicesList, ok := record.Get("choices")
	if !ok {
		return nil, fmt.Errorf("choices not found in record")
	}

	choicesData, ok := choicesList.([]any)
	if !ok {
		// Handle case where choices is null (no choices yet)
		choicesData = []any{}
	}

	edges := make([]domain.GraphEdge, 0, len(choicesData))
	for _, choiceItem := range choicesData {
		if choiceItem == nil {
			continue
		}

		choiceMap, ok := choiceItem.(map[string]any)
		if !ok {
			continue
		}

		fromID := getString(choiceMap, "from")
		toID := getString(choiceMap, "to")

		if fromID == "" || toID == "" {
			continue
		}

		edge := domain.GraphEdge{
			ID:     fmt.Sprintf("%s->%s", fromID, toID),
			Source: fromID,
			Target: toID,
			Label:  getString(choiceMap, "choice_text"),
		}

		edges = append(edges, edge)
	}

	return &domain.GraphResponse{
		Story: story,
		Nodes: nodes,
		Edges: edges,
	}, nil
}
