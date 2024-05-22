package ollama

// ChatRequestBuilder represents the chat API request.
type ChatRequestBuilder struct {
	Model      *string                        `json:"model"`
	Stream     *bool                          `json:"stream"`
	StreamFunc func(r *GCResponse, err error) `json:"-"`

	Format   *string   `json:"format"`
	Images   []string  `json:"images"`
	Raw      *bool     `json:"raw"`
	Messages []Message `json:"message"`

	Options *Options `json:"options"`
}

// WithModel sets the model used for this request.
//
// Parameters:
//   - v: The model name.
func (c ChatFunc) WithModel(v string) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Model = &v
	}
}

// WithStream passes a function to allow reading stream
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - f: The function to handle streaming
func (c ChatFunc) WithStream(v bool, f func(r *GCResponse, err error)) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Stream = &v
		r.StreamFunc = f
	}
}

// WithFormat sets the format to return a response in. Currently, the only accepted value is "json".
//
// Parameters:
//   - v: The format string.
func (c ChatFunc) WithFormat(v string) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Format = &v
	}
}

// WithImage appends an image to the message sent to  The image must be base64 encoded.
//
// Parameters:
//   - v: The base64 encoded image string.
func (c ChatFunc) WithImage(v string) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Images = append(r.Images, v)
	}
}

// WithRaw bypasses the templating system and provides a full prompt.
//
// Parameters:
//   - v: A boolean indicating whether to use raw mode.
func (c ChatFunc) WithRaw(v bool) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Raw = &v
	}
}

// WithTemperature sets the temperature for this request.
//
// Parameters:
//   - v: The temperature value.
func (c ChatFunc) WithTemperature(v float64) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
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
func (c ChatFunc) WithSeed(v int) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		if r.Options == nil {
			r.Options = &Options{}
		}

		r.Options.Seed = &v
	}
}

// WithMessage appends a new message to the request.
//
// Parameters:
//   - v: The message to append.
func (c ChatFunc) WithMessage(v Message) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		if v.Role == nil {
			v.Role = pointer("user")
		}

		r.Messages = append(r.Messages, v)
	}
}

// WithOptions sets the options for this request. It will override any settings set before, such as temperature and seed.
//
// Parameters:
//   - v: The options to set.
func (c ChatFunc) WithOptions(v Options) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Options = &v
	}
}
