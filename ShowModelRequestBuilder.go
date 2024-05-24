package ollama

// ShowModelRequestBuilder represents the model creation API request.
type ShowModelRequestBuilder struct {
	Model    *string  `json:"model"`
	System   *string  `json:"path"`
	Template *string  `json:"modelfile"`
	Options  *Options `json:"options"`
}

// WithModel sets the new model's name for this request.
//
// Parameters:
//   - v: The model name.
func (f *ShowModelInfoFunc) WithModel(v string) func(*ShowModelRequestBuilder) {
	return func(r *ShowModelRequestBuilder) {
		r.Model = &v
	}
}

// WithTemplate sets the template for this request.
//
// Parameters:
//   - v: The template string.
func (f *ShowModelInfoFunc) WithTemplate(v string) func(*ShowModelRequestBuilder) {
	return func(r *ShowModelRequestBuilder) {
		r.Template = &v
	}
}

// WithSystem sets the system for this request.
//
// Parameters:
//   - v: The system message string.
func (f *ShowModelInfoFunc) WithSystem(v string) func(*ShowModelRequestBuilder) {
	return func(r *ShowModelRequestBuilder) {
		r.System = &v
	}
}

// WithOptions sets the options for this request. It will override any settings set before, such as temperature and seed.
//
// Parameters:
//   - v: The options to set.
func (f *ShowModelInfoFunc) WithOptions(v Options) func(*ShowModelRequestBuilder) {
	return func(r *ShowModelRequestBuilder) {
		r.Options = &v
	}
}
