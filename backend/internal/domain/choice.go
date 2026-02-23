package domain

// Choice represents a choice (edge) between two scenes
type Choice struct {
	FromSceneID string `json:"from_scene_id"`
	ToSceneID   string `json:"to_scene_id"`
	ChoiceText  string `json:"choice_text"`
}

// CreateChoiceRequest represents the request to create or update a choice
type CreateChoiceRequest struct {
	FromSceneID string `json:"from_scene_id" binding:"required"`
	ToSceneID   string `json:"to_scene_id" binding:"required"`
	ChoiceText  string `json:"choice_text" binding:"required,max=200"`
}

// UpdateChoiceRequest represents the request to update choice text
type UpdateChoiceRequest struct {
	FromSceneID string `json:"from_scene_id" binding:"required"`
	ToSceneID   string `json:"to_scene_id" binding:"required"`
	ChoiceText  string `json:"choice_text" binding:"required,max=200"`
}

// DeleteChoiceRequest represents the request to delete a choice
type DeleteChoiceRequest struct {
	FromSceneID string `json:"from_scene_id" binding:"required"`
	ToSceneID   string `json:"to_scene_id" binding:"required"`
}
