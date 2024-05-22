package ollama

// PullModelRequest represents the pull model API request.
type PullModelRequest struct {
	Name       *string `json:"name"`
	Insecure   *bool   `json:"insecure"`
	Stream     *bool   `json:"stream"`
	StreamFunc func(r *StatusResponse, err error)
}

// WithName sets the model used for this request.
//
// Parameters:
//   - v: The model name.
func (f *PullModelFunc) WithName(v string) func(*PullModelRequest) {
	return func(r *PullModelRequest) {
		r.Name = &v
	}
}

// WithInsecure allows insecure connections to the library. Only use this if you are pulling from your own library during development.
//
// Parameters:
//   - v: A boolean indicating whether to insecure mode.
func (f PullModelFunc) WithInsecure(v bool) func(*PullModelRequest) {
	return func(r *PullModelRequest) {
		r.Insecure = &v
	}
}

// WithStream passes a function to allow reading stream types.
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - fc: The function to handle streaming types.
func (f *PullModelFunc) WithStream(v bool, fc func(r *StatusResponse, err error)) func(*PullModelRequest) {
	return func(r *PullModelRequest) {
		r.Stream = &v
		r.StreamFunc = fc
	}
}
