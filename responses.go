package ollama

// GCResponse represents the API response for "chat" and "generate" endpoints.
type GCResponse struct {
	Model              string  `json:"model"`
	CreatedAt          string  `json:"created_at"`
	Response           string  `json:"response"`
	Message            Message `json:"message"`
	Done               bool    `json:"done"`
	Context            []int   `json:"context"`
	TotalDuration      int64   `json:"total_duration"`
	LoadDuration       int64   `json:"load_duration"`
	PromptEvalCount    int64   `json:"prompt_eval_count"`
	PromptEvalDuration int64   `json:"prompt_eval_duration"`
	EvalCount          int64   `json:"eval_count"`
	EvalDuration       int64   `json:"eval_duration"`
}

// GenerateEmbeddingsResponse represents the API response for "generate embeddings" endpoint.
type GenerateEmbeddingsResponse struct {
	Embedding []float64 `json:"embedding"`
}

// ListLocalModelsResponse represents the response for listing local models.
type ListLocalModelsResponse struct {
	Models []Model `json:"models"`
}

// PushModelResponse represents the API response for "model push" endpoint.
type PushModelResponse struct {
	Status string `json:"status"`
	Digest string `json:"digest"`
	Total  int64  `json:"total"`
}

// ShowModelInfoResponse represents the response for showing model information.
type ShowModelInfoResponse struct {
	Modelfile  string       `json:"modelfile"`
	Parameters string       `json:"parameters"`
	Template   string       `json:"template"`
	Details    ModelDetails `json:"details"`
}

// StatusResponse represents the API response for endpoint that return status updates.
type StatusResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
