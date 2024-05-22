package ollama

import (
	"bytes"
	json2 "encoding/json"
	"io"
)

// ChatFunc performs a request to the Ollama API with the provided instructions.
// If chatId is set, it will append the messages from previous requests to the current request.
// If chatId is not found, a new chat will be generated.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type ChatFunc func(chatId *string, builder ...func(reqBuilder *ChatRequestBuilder)) (*GCResponse, error)

// GenerateFunc performs a request to the Ollama API with the provided instructions.
// If the prompt is not set, the model will be loaded into memory.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type GenerateFunc func(builder ...func(reqBuilder *GenerateRequestBuilder)) (*GCResponse, error)

// BlobCreateFunc performs a request to the Ollama API to create a new blob with the provided blob file.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type BlobCreateFunc func(digest string, data []byte) (bool, error)

// BlobCheckFunc performs a request to the Ollama API to check if a blob file exists.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type BlobCheckFunc func(digest string) (bool, error)

// CreateModelFunc performs a request to the Ollama API to create a new model with the provided model file.
// Canceled pulls are resumed from where they left off, and multiple calls will share the same download progress.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type CreateModelFunc func(builder ...func(modelFileBuilder *ModelFileBuilder)) (*StatusResponse, error)

// ListLocalModelsFunc performs a request to the Ollama API to retrieve the local models.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type ListLocalModelsFunc func() (*ListLocalModelsResponse, error)

// ShowModelInfoFunc performs a request to the Ollama API to retrieve the information of a model.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type ShowModelInfoFunc func(name string) (*ShowModelInfoResponse, error)

// CopyModelFunc performs a request to the Ollama API to copy an existing model under a different name.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type CopyModelFunc func(source, destination string) (bool, error)

// DeleteModelFunc performs a request to the Ollama API to delete a model.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type DeleteModelFunc func(name string) (bool, error)

// PullModelFunc performs a request to the Ollama API to pull model from the ollama library.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type PullModelFunc func(...func(modelFileBuilder *PullModelRequest)) (*StatusResponse, error)

// PushModelFunc performs a request to the Ollama API to push model to the ollama library.
// Requires registering for ollama.ai and adding a public key first
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type PushModelFunc func(...func(modelFileBuilder *PushModelRequest)) (*PushModelResponse, error)

// GenerateEmbeddingsFunc performs a request to the Ollama API to generate embeddings from a model.
//
// For more information about the request, see the API documentation:
// https://github.com/ollama/ollama/blob/main/docs/api.md
type GenerateEmbeddingsFunc func(...func(modelFileBuilder *GenerateEmbeddingsBuilder)) (*GenerateEmbeddingsResponse, error)

func (o *Ollama) newChatFunc() ChatFunc {
	return func(chatId *string, builder ...func(reqBuilder *ChatRequestBuilder)) (*GCResponse, error) {
		req := ChatRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.Model == nil {
			req.Model = pointer(o.defaultModel)
		}

		if req.Stream == nil {
			req.Stream = pointer(false)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(1024)
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
				req.StreamFunc(bodyTo[GCResponse](b))
			}
		}

		body, err := o.Do("/api/chat", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]GCResponse, 0)
		for _, b := range body {
			r, err := bodyTo[GCResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		// Connect types
		final := &GCResponse{}
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

			final.Message.Content = pointer(*final.Message.Content + *r.Message.Content)
			final.Message.Images = append(final.Message.Images, r.Message.Images...)

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
	return func(builder ...func(reqBuilder *GenerateRequestBuilder)) (*GCResponse, error) {
		req := GenerateRequestBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.Model == nil {
			req.Model = pointer(o.defaultModel)
		}

		if req.Stream == nil {
			req.Stream = pointer(false)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(1024)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[GCResponse](b))
			}
		}

		body, err := o.Do("/api/generate", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]GCResponse, 0)
		for _, b := range body {
			r, err := bodyTo[GCResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		// Connect types
		final := &GCResponse{}
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
	return func(digest string, data []byte) (bool, error) {
		res, err := o.Request("POST", "/api/blobs/"+digest, bytes.NewBuffer(data))
		if err != nil {
			return false, err
		}
		defer res.Body.Close()

		return res.StatusCode == 200, nil
	}
}

func (o *Ollama) newBlobCheckFunc() BlobCheckFunc {
	return func(digest string) (bool, error) {
		res, err := o.Request("HEAD", "/api/blobs/"+digest, nil)
		if err != nil {
			return false, err
		}
		defer res.Body.Close()

		return res.StatusCode == 200, nil
	}
}

func (o *Ollama) newCreateModelFunc() CreateModelFunc {
	return func(builder ...func(modelFileBuilder *ModelFileBuilder)) (*StatusResponse, error) {
		req := ModelFileBuilder{}
		for _, f := range builder {
			f(&req)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(1024)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[StatusResponse](b))
			}
		}

		body, err := o.Do("/api/create", map[string]any{
			"name":      req.Name,
			"modelfile": req.Build(o.defaultModel),
			"stream":    req.Stream,
		}, *req.StreamBufferSize, stream)
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
		res, err := o.Request("GET", "/api/tags", nil)
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
	return func(name string) (*ShowModelInfoResponse, error) {
		json, err := json2.Marshal(map[string]string{
			"name": name,
		})
		if err != nil {
			return nil, err
		}

		res, err := o.Request("POST", "/api/show", bytes.NewBuffer(json))
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
	return func(source, destination string) (bool, error) {
		json, err := json2.Marshal(map[string]string{
			"source":      source,
			"destination": destination,
		})
		if err != nil {
			return false, err
		}

		res, err := o.Request("POST", "/api/copy", bytes.NewBuffer(json))
		if err != nil {
			return false, err
		}
		defer res.Body.Close()

		return res.StatusCode == 200, nil
	}
}

func (o *Ollama) newDeleteModelFunc() DeleteModelFunc {
	return func(name string) (bool, error) {
		json, err := json2.Marshal(map[string]string{
			"name": name,
		})
		if err != nil {
			return false, err
		}

		res, err := o.Request("DELETE", "/api/delete", bytes.NewBuffer(json))
		if err != nil {
			return false, err
		}
		defer res.Body.Close()

		return res.StatusCode == 200, nil
	}
}

func (o *Ollama) newPullModelFunc() PullModelFunc {
	return func(builder ...func(modelFileBuilder *PullModelRequest)) (*StatusResponse, error) {
		req := PullModelRequest{}
		for _, f := range builder {
			f(&req)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(1024)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[StatusResponse](b))
			}
		}

		body, err := o.Do("/api/pull", req, *req.StreamBufferSize, stream)
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

func (o *Ollama) newPushModelFunc() PushModelFunc {
	return func(builder ...func(modelFileBuilder *PushModelRequest)) (*PushModelResponse, error) {
		req := PushModelRequest{}
		for _, f := range builder {
			f(&req)
		}

		if req.StreamBufferSize == nil {
			req.StreamBufferSize = pointer(1024)
		}

		var stream func(b []byte)
		if req.StreamFunc != nil {
			stream = func(b []byte) {
				req.StreamFunc(bodyTo[PushModelResponse](b))
			}
		}

		body, err := o.Do("/api/push", req, *req.StreamBufferSize, stream)
		if err != nil {
			return nil, err
		}

		resp := make([]PushModelResponse, 0)
		for _, b := range body {
			r, err := bodyTo[PushModelResponse](b)
			if err != nil {
				return nil, err
			}
			resp = append(resp, *r)
		}

		final := &PushModelResponse{}
		for _, r := range resp {
			final.Status += r.Status + "\n"
		}

		return final, nil
	}
}

func (o *Ollama) newGenerateEmbeddingsFunc() GenerateEmbeddingsFunc {
	return func(builder ...func(modelFileBuilder *GenerateEmbeddingsBuilder)) (*GenerateEmbeddingsResponse, error) {
		req := GenerateEmbeddingsBuilder{}
		for _, f := range builder {
			f(&req)
		}

		body, err := o.Do("/api/embeddings", req, 0, nil)
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
