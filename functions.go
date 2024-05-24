package ollama

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatFunc performs a request to the Ollama API with the provided instructions.
// If chatId is set, it will append the messages from previous requests to the current request.
// If chatId is not found, a new chat will be generated.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type ChatFunc func(chatId *string, builder ...func(reqBuilder *ChatRequestBuilder)) (*ChatResponse, error)

// GenerateFunc performs a request to the Ollama API with the provided instructions.
// If the prompt is not set, the model will be loaded into memory.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type GenerateFunc func(builder ...func(reqBuilder *GenerateRequestBuilder)) (*GenerateResponse, error)

// BlobCreateFunc performs a request to the Ollama API to create a new blob with the provided blob file.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type BlobCreateFunc func(digest string, data []byte) error

// BlobCheckFunc performs a request to the Ollama API to check if a blob file exists.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type BlobCheckFunc func(digest string) error

// CreateModelFunc performs a request to the Ollama API to create a new model with the provided model file.
// Canceled pulls are resumed from where they left off, and multiple calls will share the same download progress.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type CreateModelFunc func(builder ...func(modelFileBuilder *ModelFileRequestBuilder)) (*StatusResponse, error)

// ListLocalModelsFunc performs a request to the Ollama API to retrieve the local models.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type ListLocalModelsFunc func() (*ListLocalModelsResponse, error)

// ShowModelInfoFunc performs a request to the Ollama API to retrieve the information of a model.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type ShowModelInfoFunc func(builder ...func(reqBuilder *ShowModelRequestBuilder)) (*ShowModelInfoResponse, error)

// CopyModelFunc performs a request to the Ollama API to copy an existing model under a different name.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type CopyModelFunc func(source, destination string) error

// DeleteModelFunc performs a request to the Ollama API to delete a model.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type DeleteModelFunc func(name string) error

// PullModelFunc performs a request to the Ollama API to pull model from the ollama library.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type PullModelFunc func(...func(modelFileBuilder *PullModelRequestBuilder)) (*PushPullModelResponse, error)

// PushModelFunc performs a request to the Ollama API to push model to the ollama library.
// Requires registering for ollama.ai and adding a public key first
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type PushModelFunc func(...func(modelFileBuilder *PushModelRequestBuilder)) (*PushPullModelResponse, error)

// GenerateEmbeddingsFunc performs a request to the Ollama API to generate embeddings from a model.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type GenerateEmbeddingsFunc func(...func(modelFileBuilder *GenerateEmbeddingsRequestBuilder)) (*GenerateEmbeddingsResponse, error)

// VersionFunc performs a request to the Ollama API and returns the Ollama server version as a string.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type VersionFunc func() (*VersionResponse, error)

func (o *Ollama) newChatFunc() ChatFunc {
	return func(chatId *string, builder ...func(reqBuilder *ChatRequestBuilder)) (*ChatResponse, error) {
		req := ChatRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.Stream == nil {
			req.Stream = pointer(false)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(512000)
		}

		// Include chat history or create a new chat
		if chatId != nil {
			chat := o.chats[*chatId]
			if chat == nil {
				chat = &Chat{
					ID:       *chatId,
					Messages: make([]Message, 0),
				}
				o.chats[*chatId] = chat
			}

			for _, chat := range chat.Messages {
				req.Messages = append([]Message{chat}, req.Messages...)
			}
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[ChatResponse](b))
			}
		}

		body, err := o.stream(http.MethodPost, "/api/chat", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]ChatResponse, 0)
		for _, b := range body {
			r, err := bodyTo[ChatResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		// Connect types
		final := &ChatResponse{}
		for i, r := range resp {
			if i == 0 {
				final.Model = r.Model
				final.CreatedAt = r.CreatedAt
				final.Message = Message{
					Role:    r.Message.Role,
					Content: pointer(""),
				}
				final.Done = r.Done
			}

			if r.Message.Content != nil {
				final.Message.Content = pointer(*final.Message.Content + *r.Message.Content)
			}

			if r.Message.Images != nil && len(r.Message.Images) > 0 {
				final.Message.Images = append(final.Message.Images, r.Message.Images...)
			}

			if i == len(resp)-1 {
				final.TotalDuration = r.TotalDuration
				final.LoadDuration = r.LoadDuration
				final.PromptEvalCount = r.PromptEvalCount
				final.PromptEvalDuration = r.PromptEvalDuration
				final.EvalCount = r.EvalCount
				final.EvalDuration = r.EvalDuration
				final.Context = r.Context
			}
		}

		if chatId != nil {
			o.chats[*chatId].AddMessage(final.Message)
		}

		return final, nil
	}
}

func (o *Ollama) newGenerateFunc() GenerateFunc {
	return func(builder ...func(reqBuilder *GenerateRequestBuilder)) (*GenerateResponse, error) {
		req := GenerateRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.Stream == nil {
			req.Stream = pointer(false)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(512000)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[GenerateResponse](b))
			}
		}

		body, err := o.stream(http.MethodPost, "/api/generate", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]GenerateResponse, 0)
		for _, b := range body {
			r, err := bodyTo[GenerateResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		// Connect types
		final := &GenerateResponse{}
		for i, r := range resp {
			if i == 0 {
				final.Model = r.Model
				final.CreatedAt = r.CreatedAt
				final.Done = r.Done
			}

			final.Response += r.Response

			if i == len(resp)-1 {
				final.TotalDuration = r.TotalDuration
				final.LoadDuration = r.LoadDuration
				final.PromptEvalCount = r.PromptEvalCount
				final.PromptEvalDuration = r.PromptEvalDuration
				final.EvalCount = r.EvalCount
				final.EvalDuration = r.EvalDuration
				final.Context = r.Context
			}
		}

		return final, nil
	}
}

func (o *Ollama) newBlobCreateFunc() BlobCreateFunc {
	return func(digest string, data []byte) error {
		res, err := o.request(http.MethodPost, "/api/blobs/"+digest, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		return nil
	}
}

func (o *Ollama) newBlobCheckFunc() BlobCheckFunc {
	return func(digest string) error {
		res, err := o.request(http.MethodHead, "/api/blobs/"+digest, nil)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		return nil
	}
}

func (o *Ollama) newCreateModelFunc() CreateModelFunc {
	return func(builder ...func(modelFileBuilder *ModelFileRequestBuilder)) (*StatusResponse, error) {
		req := ModelFileRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(512000)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[StatusResponse](b))
			}
		}

		req.Modelfile = pointer(req.Build())

		body, err := o.stream(http.MethodPost, "/api/create", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]StatusResponse, 0)
		for _, b := range body {
			r, err := bodyTo[StatusResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		final := &StatusResponse{}
		for _, r := range resp {
			final.Status += r.Status + "\n"
		}

		return final, nil
	}
}

func (o *Ollama) newListLocalModelsFunc() ListLocalModelsFunc {
	return func() (*ListLocalModelsResponse, error) {
		res, err := o.request(http.MethodGet, "/api/tags", nil)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return bodyTo[ListLocalModelsResponse](body)
	}
}

func (o *Ollama) newShowModelInfoFunc() ShowModelInfoFunc {
	return func(builder ...func(reqBuilder *ShowModelRequestBuilder)) (*ShowModelInfoResponse, error) {
		req := ShowModelRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		json, err := json2.Marshal(req)
		if err != nil {
			return nil, err
		}

		res, err := o.request(http.MethodPost, "/api/show", bytes.NewBuffer(json))
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return bodyTo[ShowModelInfoResponse](body)
	}
}

func (o *Ollama) newCopyModelFunc() CopyModelFunc {
	return func(source, destination string) error {
		json, err := json2.Marshal(map[string]string{
			"source":      source,
			"destination": destination,
		})
		if err != nil {
			return err
		}

		res, err := o.request(http.MethodPost, "/api/copy", bytes.NewBuffer(json))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		return nil
	}
}

func (o *Ollama) newDeleteModelFunc() DeleteModelFunc {
	return func(model string) error {
		json, err := json2.Marshal(map[string]string{
			"model": model,
		})
		if err != nil {
			return err
		}

		res, err := o.request(http.MethodDelete, "/api/delete", bytes.NewBuffer(json))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		return nil
	}
}

func (o *Ollama) newPullModelFunc() PullModelFunc {
	return func(builder ...func(modelFileBuilder *PullModelRequestBuilder)) (*PushPullModelResponse, error) {
		req := PullModelRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(512000)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[PushPullModelResponse](b))
			}
		}

		body, err := o.stream(http.MethodPost, "/api/pull", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]PushPullModelResponse, 0)
		for _, b := range body {
			r, err := bodyTo[PushPullModelResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		final := &PushPullModelResponse{}
		for _, r := range resp {
			if len(r.Status) != 0 {
				final.Status += r.Status + "\n"
			}

			if len(r.Error) != 0 {
				final.Error += r.Error + "\n"
			}
		}

		return final, nil
	}
}

func (o *Ollama) newPushModelFunc() PushModelFunc {
	return func(builder ...func(modelFileBuilder *PushModelRequestBuilder)) (*PushPullModelResponse, error) {
		req := PushModelRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(512000)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[PushPullModelResponse](b))
			}
		}

		body, err := o.stream(http.MethodPost, "/api/push", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]PushPullModelResponse, 0)
		for _, b := range body {
			r, err := bodyTo[PushPullModelResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		final := &PushPullModelResponse{}
		for _, r := range resp {
			final.Status += r.Status + "\n"
		}

		return final, nil
	}
}

func (o *Ollama) newGenerateEmbeddingsFunc() GenerateEmbeddingsFunc {
	return func(builder ...func(modelFileBuilder *GenerateEmbeddingsRequestBuilder)) (*GenerateEmbeddingsResponse, error) {
		req := GenerateEmbeddingsRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		body, err := o.stream(http.MethodPost, "/api/embeddings", req, 0, nil)
		if err != nil {
			return nil, err
		}

		r, err := bodyTo[GenerateEmbeddingsResponse](body[0])
		if err != nil {
			return nil, err
		}

		return r, nil
	}
}

func (o *Ollama) newVersionFunc() VersionFunc {
	return func() (*VersionResponse, error) {
		res, err := o.request(http.MethodGet, "/api/version", nil)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		respBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("status code: %d, failed to read response body: %w", res.StatusCode, err)
		}

		r, err := bodyTo[VersionResponse](respBody)
		if err != nil {
			return nil, err
		}

		return r, nil
	}
}
