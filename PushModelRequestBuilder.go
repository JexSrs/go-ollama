package ollama

// PushModelRequestBuilder represents the push model API request.
type PushModelRequestBuilder struct {
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
func (f *PushModelFunc) WithModel(v string) func(*PushModelRequestBuilder) {
	return func(r *PushModelRequestBuilder) {
		r.Model = &v
	}
}

// WithInsecure allows insecure connections to the library. Only use this if you are pulling from your own library during development.
//
// Parameters:
//   - v: A boolean indicating whether to insecure mode.
func (f *PushModelFunc) WithInsecure(v bool) func(*PushModelRequestBuilder) {
	return func(r *PushModelRequestBuilder) {
		r.Insecure = &v
	}
}

// WithUsername sets the username used for this request.
//
// Parameters:
//   - v: The username.
func (f *PushModelFunc) WithUsername(v string) func(*PushModelRequestBuilder) {
	return func(r *PushModelRequestBuilder) {
		r.Username = &v
	}
}

// WithPassword sets the password used for this request.
//
// Parameters:
//   - v: The password.
func (f *PushModelFunc) WithPassword(v string) func(*PushModelRequestBuilder) {
	return func(r *PushModelRequestBuilder) {
		r.Password = &v
	}
}

// WithStream passes a function to allow reading stream types.
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - bufferSize: The size of the streamed buffer
//   - fc: The function to handle streaming types.
func (f *PushModelFunc) WithStream(v bool, bufferSize int, fc func(r *PushPullModelResponse, err error)) func(*PushModelRequestBuilder) {
	return func(r *PushModelRequestBuilder) {
		r.Stream = &v
		r.StreamBufferSize = &bufferSize
		r.StreamFunc = fc
	}
}
