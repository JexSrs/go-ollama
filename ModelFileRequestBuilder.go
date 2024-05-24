package ollama

// ModelFileRequestBuilder represents the model creation API request.
type ModelFileRequestBuilder struct {
	Model     *string `json:"model"`
	Path      *string `json:"path"`
	Modelfile *string `json:"modelfile"`
	Quantize  *string `json:"quantize"`

	Stream           *bool                              `json:"stream"`
	StreamBufferSize *int                               `json:"-"`
	StreamFunc       func(r *StatusResponse, err error) `json:"-"`

	from       *string
	parameters []Parameter
	template   *string
	system     *string
	adapter    *string
	license    *string

	messages []Message
}

// Parameter represents a parameter sent to the API,
type Parameter struct {
	Key   string
	Value string
}

// WithModel sets the new model's name for this request.
//
// Parameters:
//   - v: The model name.
func (f *CreateModelFunc) WithModel(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.Model = &v
	}
}

// WithPath sets the path for this request.
//
// Parameters:
//   - v: The path.
func (f *CreateModelFunc) WithPath(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.Path = &v
	}
}

// WithStream passes a function to allow reading stream types.
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - bufferSize: The size of the streamed buffer
//   - fc: The function to handle streaming types.
func (f *CreateModelFunc) WithStream(v bool, bufferSize int, fc func(r *StatusResponse, err error)) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.Stream = &v
		r.StreamBufferSize = &bufferSize
		r.StreamFunc = fc
	}
}

// WithQuantize sets the quantize for this request.
//
// Parameters:
//   - v: The quantize value.
func (f *CreateModelFunc) WithQuantize(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.Quantize = &v
	}
}

// WithFrom defines the base model to use.
//
// Parameters:
//   - v: The base model string.
func (f *CreateModelFunc) WithFrom(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.from = &v
	}
}

// WithTemplate sets the full prompt template to be sent to the model.
//
// Parameters:
//   - v: The template string.
func (f *CreateModelFunc) WithTemplate(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.template = &v
	}
}

// WithSystem specifies the system message that will be set in the template.
//
// Parameters:
//   - v: The system message string.
func (f *CreateModelFunc) WithSystem(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.system = &v
	}
}

// WithParameter appends a new parameter for how Ollama will run the model.
//
// Parameters:
//   - v: The parameter to append.
func (f *CreateModelFunc) WithParameter(v Parameter) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.parameters = append(r.parameters, v)
	}
}

// WithAdapter defines the (Q)LoRA adapters to apply to the model.
//
// Parameters:
//   - v: The adapter string.
func (f *CreateModelFunc) WithAdapter(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.adapter = &v
	}
}

// WithLicense specifies the legal license.
//
// Parameters:
//   - v: The license string.
func (f *CreateModelFunc) WithLicense(v string) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.license = &v
	}
}

// WithMessage appends a new message to the message history.
//
// Parameters:
//   - v: The message to append.
func (f *CreateModelFunc) WithMessage(v Message) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.messages = append(r.messages, v)
	}
}

// WithChat appends all the messages from a chat to the message history.
//
// Parameters:
//   - chat: The chat whose messages to append.
func (f *CreateModelFunc) WithChat(chat *Chat) func(*ModelFileRequestBuilder) {
	return func(r *ModelFileRequestBuilder) {
		r.messages = append(r.messages, chat.Messages...)
	}
}

// Build generates the ModelFile.
//
// Parameters:
//   - defaultModel: The default model string.
func (m *ModelFileRequestBuilder) Build() string {
	r := ""

	if m.from != nil {
		r += "FROM " + *m.from + "\n"
	}

	for _, p := range m.parameters {
		r += "PARAMETER " + p.Key + " " + p.Value + "\n"
	}

	if m.template != nil {
		r += "TEMPLATE \"\"\"" + *m.template + "\"\"\"\n"
	}

	if m.system != nil {
		r += "SYSTEM \"\"\"" + *m.system + "\"\"\"\n"
	}

	if m.adapter != nil {
		r += "ADAPTER " + *m.adapter + "\n"
	}

	if m.license != nil {
		r += "LICENSE\n\"\"\"" + *m.license + "\n\"\"\"\n"
	}

	for _, p := range m.messages {
		r += "MESSAGE " + *p.Role + " " + *p.Content + "\n"
	}

	return r
}
