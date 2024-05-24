package ollama

// PullModelRequestBuilder represents the pull model API request.
type PullModelRequestBuilder struct {
	Model    *string `json:"model"`
	Insecure *bool   `json:"insecure"`
	Username *string `json:"username"`
	Password *string `json:"password"`

	Stream           *bool                                     `json:"stream"`
	StreamBufferSize *int                                      `json:"-"`
	StreamFunc       func(r *PushPullModelResponse, err error) `json:"-"`
}

// WithModel sets the model used for this request.
//
// Parameters:
//   - v: The model name.
func (f *PullModelFunc) WithModel(v string) func(*PullModelRequestBuilder) {
	return func(r *PullModelRequestBuilder) {
		r.Model = &v
	}
}

// WithInsecure allows insecure connections to the library. Only use this if you are pulling from your own library during development.
//
// Parameters:
//   - v: A boolean indicating whether to insecure mode.
func (f *PullModelFunc) WithInsecure(v bool) func(*PullModelRequestBuilder) {
	return func(r *PullModelRequestBuilder) {
		r.Insecure = &v
	}
}

// WithUsername sets the username used for this request.
//
// Parameters:
//   - v: The username.
func (f *PullModelFunc) WithUsername(v string) func(*PullModelRequestBuilder) {
	return func(r *PullModelRequestBuilder) {
		r.Username = &v
	}
}

// WithPassword sets the password used for this request.
//
// Parameters:
//   - v: The password.
func (f *PullModelFunc) WithPassword(v string) func(*PullModelRequestBuilder) {
	return func(r *PullModelRequestBuilder) {
		r.Password = &v
	}
}

// WithStream passes a function to allow reading stream types.
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - bufferSize: The size of the streamed buffer
//   - fc: The function to handle streaming types.
func (f *PullModelFunc) WithStream(v bool, bufferSize int, fc func(r *PushPullModelResponse, err error)) func(*PullModelRequestBuilder) {
	return func(r *PullModelRequestBuilder) {
		r.Stream = &v
		r.StreamBufferSize = &bufferSize
		r.StreamFunc = fc
	}
}
