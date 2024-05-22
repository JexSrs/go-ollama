package ollama

// GenerateRequestBuilder represents the generate API request.
type GenerateRequestBuilder struct {
	Model      *string                        `json:"model"`
	Prompt     *string                        `json:"prompt"`
	Stream     *bool                          `json:"stream"`
	StreamFunc func(r *GCResponse, err error) `json:"-"`

	// Format the format to return a response in. Currently, the only accepted value is json
	Format *string  `json:"format"`
	Images []string `json:"images"`
	Raw    *bool    `json:"raw"`

	Options *Options `json:"options"`
}

// WithModel sets the model used for this request.
//
// Parameters:
//   - v: The model name.
func (c GenerateFunc) WithModel(v string) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		r.Model = &v
	}
}

// WithPrompt sets the prompt for this request.
//
// Parameters:
//   - v: The prompt string.
func (c GenerateFunc) WithPrompt(v string) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		r.Prompt = &v
	}
}

// WithStream passes a function to allow reading stream types.
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - f: The function to handle streaming types.
func (c GenerateFunc) WithStream(v bool, f func(r *GCResponse, err error)) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		r.Stream = &v
		r.StreamFunc = f
	}
}

// WithFormat sets the format to return a response in. Currently, the only accepted value is "json".
//
// Parameters:
//   - v: The format string.
func (c GenerateFunc) WithFormat(v string) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		r.Format = &v
	}
}

// WithImage appends an image to the message sent to Ollama. The image must be base64 encoded.
//
// Parameters:
//   - v: The base64 encoded image string.
func (c GenerateFunc) WithImage(v string) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		r.Images = append(r.Images, v)
	}
}

// WithRaw bypasses the templating system and provides a full prompt.
//
// Parameters:
//   - v: A boolean indicating whether to use raw mode.
func (c GenerateFunc) WithRaw(v bool) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		r.Raw = &v
	}
}

// WithTemperature sets the temperature for this request.
//
// Parameters:
//   - v: The temperature value.
func (c GenerateFunc) WithTemperature(v float64) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		if r.Options == nil {
			r.Options = &Options{}
		}

		r.Options.Temperature = &v
	}
}

// WithSeed sets the seed for this request.
//
// Parameters:
//   - v: The seed value.
func (c GenerateFunc) WithSeed(v int) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		if r.Options == nil {
			r.Options = &Options{}
		}

		r.Options.Seed = &v
	}
}

// WithOptions sets the options for this request. It will override any settings set before, such as temperature and seed.
//
// Parameters:
//   - v: The options to set.
func (c GenerateFunc) WithOptions(v Options) func(*GenerateRequestBuilder) {
	return func(r *GenerateRequestBuilder) {
		r.Options = &v
	}
}
