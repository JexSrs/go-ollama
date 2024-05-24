package ollama

import "time"

// GenerateResponse represents the API response for "generate" endpoint.
type GenerateResponse struct {
	Model      string `json:"model"`       // Is the model name that generated the response.
	CreatedAt  string `json:"created_at"`  // Is the timestamp of the response.
	Response   string `json:"response"`    // Is the textual response itself.
	Done       bool   `json:"done"`        // Specifies if the response is complete.
	DoneReason string `json:"done_reason"` // The reason the model stopped generating text.
	Context    []int  `json:"context"`     // Is an encoding of the conversation used in this response; this can be sent in the next request to keep a conversational memory.

	Metrics
}

// ChatResponse represents the API response for "chat" endpoint.
type ChatResponse struct {
	Model      string  `json:"model"`      // Is the model name that generated the response.
	CreatedAt  string  `json:"created_at"` // Is the timestamp of the response.
	Message    Message `json:"message"`
	Done       bool    `json:"done"`        // Specifies if the response is complete.
	DoneReason string  `json:"done_reason"` // The reason the model stopped generating text.
	Context    []int   `json:"context"`     // Is an encoding of the conversation used in this response; this can be sent in the next request to keep a conversational memory.

	Metrics
}

// GenerateEmbeddingsResponse represents the API response for "generate embeddings" endpoint.
type GenerateEmbeddingsResponse struct {
	Embedding []float64 `json:"embedding"`
}

// ListLocalModelsResponse represents the response for listing local models.
type ListLocalModelsResponse struct {
	Models []ModelResponse `json:"models"`
}

// ModelResponse represents a model's metadata.
type ModelResponse struct {
	Name       string       `json:"name"`
	Model      string       `json:"model"`
	ModifiedAt string       `json:"modified_at"`
	Size       int64        `json:"size"`
	Digest     string       `json:"digest"`
	Details    ModelDetails `json:"details"`
	ExpiresAt  time.Time    `json:"expires_at"`
	SizeVRAM   int64        `json:"size_vram"`
}

// PushPullModelResponse represents the API response for "model push" endpoint.
type PushPullModelResponse struct {
	Status    string `json:"status"`
	Error     string `json:"error"`
	Digest    string `json:"digest"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
}

// ShowModelInfoResponse represents the response for showing model information.
type ShowModelInfoResponse struct {
	License    string       `json:"license"`
	Modelfile  string       `json:"modelfile"`
	Parameters string       `json:"parameters"`
	Template   string       `json:"template"`
	System     string       `json:"system"`
	Details    ModelDetails `json:"details"`
	Messages   []Message    `json:"messages"`
}

// StatusResponse represents the API response for endpoint that return status updates.
type StatusResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// VersionResponse represents the API response for the version endpoint.
type VersionResponse struct {
	Version string `json:"version"`
}
