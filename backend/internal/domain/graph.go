package domain

// GraphResponse represents the full story graph for Vue Flow
type GraphResponse struct {
	Story StorySummary    `json:"story"`
	Nodes []GraphNode     `json:"nodes"`
	Edges []GraphEdge     `json:"edges"`
}

// StorySummary represents minimal story info for graph response
type StorySummary struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// GraphNode represents a node in Vue Flow format
type GraphNode struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Position Position    `json:"position"`
	Data     NodeData    `json:"data"`
}

// Position represents node position in Vue Flow
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeData represents the data payload for a scene node
type NodeData struct {
	Content string  `json:"content"`
	IsStart bool    `json:"is_start"`
	IsEnd   bool    `json:"is_end"`
	ImageURL *string `json:"image_url,omitempty"`
}

// GraphEdge represents an edge in Vue Flow format
type GraphEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
}
