package ollama

// PushModelRequest represents the push model API request.
type PushModelRequest struct {
	Name       *string `json:"name"`
	Insecure   *bool   `json:"insecure"`
	Stream     *bool   `json:"stream"`
	StreamFunc func(r *PushModelResponse, err error)
}

// WithName sets the model used for this request.
//
// Parameters:
//   - v: The model name.
func (f *PushModelFunc) WithName(v string) func(*PushModelRequest) {
	return func(r *PushModelRequest) {
		r.Name = &v
	}
}

// WithInsecure allows insecure connections to the library. Only use this if you are pulling from your own library during development.
//
// Parameters:
//   - v: A boolean indicating whether to insecure mode.
func (f PushModelFunc) WithInsecure(v bool) func(*PushModelRequest) {
	return func(r *PushModelRequest) {
		r.Insecure = &v
	}
}

// WithStream passes a function to allow reading stream types.
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - fc: The function to handle streaming types.
func (f *PushModelFunc) WithStream(v bool, fc func(r *PushModelResponse, err error)) func(*PushModelRequest) {
	return func(r *PushModelRequest) {
		r.Stream = &v
		r.StreamFunc = fc
	}
}
