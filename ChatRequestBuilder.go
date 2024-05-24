package ollama

// ChatRequestBuilder represents the chat API request.
type ChatRequestBuilder struct {
	Model     *string   `json:"model"`
	Format    *string   `json:"format"`
	Raw       *bool     `json:"raw"`
	Messages  []Message `json:"messages"`
	KeepAlive *string   `json:"keep_alive,omitempty"`
	Options   *Options  `json:"options"`

	Stream           *bool                            `json:"stream"`
	StreamBufferSize *int                             `json:"-"`
	StreamFunc       func(r *ChatResponse, err error) `json:"-"`
}

// WithModel sets the model used for this request.
//
// Parameters:
//   - v: The model name.
func (f *ChatFunc) WithModel(v string) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Model = &v
	}
}

// WithStream passes a function to allow reading stream
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - bufferSize: The size of the streamed buffer
//   - f: The function to handle streaming
func (f *ChatFunc) WithStream(v bool, bufferSize int, fn func(r *ChatResponse, err error)) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Stream = &v
		r.StreamBufferSize = &bufferSize
		r.StreamFunc = fn
	}
}

// WithFormat sets the format to return a response in. Currently, the only accepted value is "json".
//
// Parameters:
//   - v: The format string.
func (f *ChatFunc) WithFormat(v string) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Format = &v
	}
}

// WithRaw bypasses the templating system and provides a full prompt.
//
// Parameters:
//   - v: A boolean indicating whether to use raw mode.
func (f *ChatFunc) WithRaw(v bool) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Raw = &v
	}
}

// WithTemperature sets the temperature for this request.
//
// Parameters:
//   - v: The temperature value.
func (f *ChatFunc) WithTemperature(v float64) func(*ChatRequestBuilder) {
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
func (f *ChatFunc) WithSeed(v int) func(*ChatRequestBuilder) {
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
func (f *ChatFunc) WithMessage(v Message) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		if v.Role == nil {
			v.Role = pointer("user")
		}

		r.Messages = append(r.Messages, v)
	}
}

// WithKeepAlive controls how long the model will stay loaded into memory following the request.
//
// Parameters:
//   - v: The keep alive duration.
func (f *ChatFunc) WithKeepAlive(v string) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.KeepAlive = &v
	}
}

// WithOptions sets the options for this request. It will override any settings set before, such as temperature and seed.
//
// Parameters:
//   - v: The options to set.
func (f *ChatFunc) WithOptions(v Options) func(*ChatRequestBuilder) {
	return func(r *ChatRequestBuilder) {
		r.Options = &v
	}
}
