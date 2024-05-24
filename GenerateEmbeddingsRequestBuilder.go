package ollama

// GenerateEmbeddingsRequestBuilder represents the generate embeddings API request.
type GenerateEmbeddingsRequestBuilder struct {
	Model     *string  `json:"model"`
	Prompt    *string  `json:"prompt"`
	KeepAlive *string  `json:"keep_alive"`
	Options   *Options `json:"options"`
}

// WithModel sets the model used for this request.
//
// Parameters:
//   - v: The model name.
func (c GenerateEmbeddingsFunc) WithModel(v string) func(*GenerateEmbeddingsRequestBuilder) {
	return func(r *GenerateEmbeddingsRequestBuilder) {
		r.Model = &v
	}
}

// WithPrompt sets the prompt for this request.
//
// Parameters:
//   - v: The prompt string.
func (c GenerateEmbeddingsFunc) WithPrompt(v string) func(*GenerateEmbeddingsRequestBuilder) {
	return func(r *GenerateEmbeddingsRequestBuilder) {
		r.Prompt = &v
	}
}

// WithKeepAlive controls how long the model will stay loaded into memory following the request (default: 5m).
//
// Parameters:
//   - v: The keep alive string.
func (c GenerateEmbeddingsFunc) WithKeepAlive(v string) func(*GenerateEmbeddingsRequestBuilder) {
	return func(r *GenerateEmbeddingsRequestBuilder) {
		r.KeepAlive = &v
	}
}

// WithOptions sets the options for this request. It will override any settings set before, such as temperature and seed.
//
// Parameters:
//   - v: The options to set.
func (c GenerateEmbeddingsFunc) WithOptions(v Options) func(*GenerateEmbeddingsRequestBuilder) {
	return func(r *GenerateEmbeddingsRequestBuilder) {
		r.Options = &v
	}
}
