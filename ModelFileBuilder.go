package ollama

// ModelFileBuilder represents the model creation API request.
type ModelFileBuilder struct {
	Name *string

	Stream     *bool
	StreamFunc func(r *StatusResponse, err error)

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
	key   string
	value string
}

// WithStream passes a function to allow reading stream types.
//
// Parameters:
//   - v: A boolean indicating whether to use streaming.
//   - fc: The function to handle streaming types.
func (f *CreateModelFunc) WithStream(v bool, fc func(r *StatusResponse, err error)) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.Stream = &v
		r.StreamFunc = fc
	}
}

// WithFrom defines the base model to use.
//
// Parameters:
//   - v: The base model string.
func (f *CreateModelFunc) WithFrom(v string) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.from = &v
	}
}

// WithTemplate sets the full prompt template to be sent to the model.
//
// Parameters:
//   - v: The template string.
func (f *CreateModelFunc) WithTemplate(v string) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.template = &v
	}
}

// WithSystem specifies the system message that will be set in the template.
//
// Parameters:
//   - v: The system message string.
func (f *CreateModelFunc) WithSystem(v string) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.system = &v
	}
}

// WithParameter appends a new parameter for how Ollama will run the model.
//
// Parameters:
//   - v: The parameter to append.
func (f *CreateModelFunc) WithParameter(v Parameter) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.parameters = append(r.parameters, v)
	}
}

// WithAdapter defines the (Q)LoRA adapters to apply to the model.
//
// Parameters:
//   - v: The adapter string.
func (f *CreateModelFunc) WithAdapter(v string) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.adapter = &v
	}
}

// WithLicense specifies the legal license.
//
// Parameters:
//   - v: The license string.
func (f *CreateModelFunc) WithLicense(v string) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.license = &v
	}
}

// WithMessage appends a new message to the message history.
//
// Parameters:
//   - v: The message to append.
func (f *CreateModelFunc) WithMessage(v Message) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.messages = append(r.messages, v)
	}
}

// WithChat appends all the messages from a chat to the message history.
//
// Parameters:
//   - chat: The chat whose messages to append.
func (f *CreateModelFunc) WithChat(chat *Chat) func(*ModelFileBuilder) {
	return func(r *ModelFileBuilder) {
		r.messages = append(r.messages, chat.Messages...)
	}
}

// Build generates the ModelFile.
//
// Parameters:
//   - defaultModel: The default model string.
func (m *ModelFileBuilder) Build(defaultModel string) string {
	r := ""

	if m.from == nil {
		m.from = pointer(defaultModel)
	}

	if m.from != nil {
		r += "FROM " + *m.from + "\n"
	}

	for _, p := range m.parameters {
		r += "PARAMETER " + p.key + " " + p.value + "\n"
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
