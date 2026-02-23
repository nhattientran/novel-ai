package domain

// Scene represents a scene (node) in a story graph
type Scene struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url"`
	IsEnd    bool    `json:"is_end"`
	PosX     float64 `json:"pos_x"`
	PosY     float64 `json:"pos_y"`
}

// CreateSceneRequest represents the request to create a new scene
type CreateSceneRequest struct {
	Content  string  `json:"content" binding:"required,max=5000"`
	ImageURL *string `json:"image_url"`
	IsEnd    bool    `json:"is_end"`
	PosX     float64 `json:"pos_x"`
	PosY     float64 `json:"pos_y"`
}

// UpdateSceneRequest represents the request to update a scene (partial update)
type UpdateSceneRequest struct {
	Content  *string  `json:"content" binding:"omitempty,max=5000"`
	ImageURL **string `json:"image_url"`
	IsEnd    *bool    `json:"is_end"`
	PosX     *float64 `json:"pos_x"`
	PosY     *float64 `json:"pos_y"`
}

// SceneWithStartFlag extends Scene with is_start flag for graph response
type SceneWithStartFlag struct {
	Scene
	IsStart bool `json:"is_start"`
}
