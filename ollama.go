package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Ollama represents a client for interacting with the Ollama API.
type Ollama struct {
	url     url.URL
	Http    *http.Client
	chats   map[string]*Chat
	headers map[string][]string

	Chat     ChatFunc
	Generate GenerateFunc

	Blobs struct {
		Check  BlobCheckFunc
		Create BlobCreateFunc
	}

	Models struct {
		Create   CreateModelFunc
		List     ListLocalModelsFunc
		ShowInfo ShowModelInfoFunc
		Copy     CopyModelFunc
		Delete   DeleteModelFunc
		Pull     PullModelFunc
		Push     PushModelFunc
	}

	GenerateEmbeddings GenerateEmbeddingsFunc
}

// New creates a new Ollama client that points to the specified URL.
// It initializes the client with default settings and available API functions.
//
// Example:
//
//	llm := New("http://api.ollama.com")
func New(v url.URL) *Ollama {
	o := &Ollama{
		url:     v,
		Http:    &http.Client{},
		chats:   make(map[string]*Chat),
		headers: make(map[string][]string),
	}

	o.Chat = o.newChatFunc()
	o.Generate = o.newGenerateFunc()

	o.Blobs.Check = o.newBlobCheckFunc()
	o.Blobs.Create = o.newBlobCreateFunc()

	o.Models.Create = o.newCreateModelFunc()
	o.Models.List = o.newListLocalModelsFunc()
	o.Models.ShowInfo = o.newShowModelInfoFunc()
	o.Models.Copy = o.newCopyModelFunc()
	o.Models.Delete = o.newDeleteModelFunc()
	o.Models.Pull = o.newPullModelFunc()
	o.Models.Push = o.newPushModelFunc()

	o.GenerateEmbeddings = o.newGenerateEmbeddingsFunc()

	return o
}

// PreloadChat preloads a chat into the client's chat map.
//
// Parameters:
//   - chat: The chat to preload.
func (o *Ollama) PreloadChat(chat Chat) {
	o.chats[chat.ID] = &chat
}

// GetChat retrieves a chat by its ID.
//
// Parameters:
//   - id: The ID of the chat.
//
// Returns:
//   - A pointer to the Chat if found, or nil if not found.
func (o *Ollama) GetChat(id string) *Chat {
	return o.chats[id]
}

// DeleteChat removes a chat by its ID.
//
// Parameters:
//   - id: The ID of the chat to remove.
func (o *Ollama) DeleteChat(id string) {
	delete(o.chats, id)
}

// DeleteAllChats removes all chats from the client's chat map.
func (o *Ollama) DeleteAllChats() {
	o.chats = make(map[string]*Chat, 0)
}

// SetHeaders sets the headers for all the requests.
func (o *Ollama) SetHeaders(key string, value []string) {
	o.headers[key] = value
}

func (o *Ollama) stream(method, path string, data interface{}, maxBufferSize int, streamFunc func(b []byte)) ([][]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := o.request(method, path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res [][]byte
	var buffer bytes.Buffer

	for {
		buf := make([]byte, maxBufferSize)
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			break
		}

		bigChunk := splitJSONObjects(buf[:n])
		for _, chunk := range bigChunk {
			res = append(res, chunk)
			buffer.Write(chunk)

			if streamFunc != nil {
				streamFunc(chunk)
			}
		}
	}

	return res, nil
}

func (o *Ollama) request(method, path string, body io.Reader) (*http.Response, error) {
	httpReq, err := http.NewRequest(method, o.url.JoinPath(path).String(), body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	for k, v := range o.headers {
		httpReq.Header.Del(k)
		for _, h := range v {
			httpReq.Header.Add(k, h)
		}
	}

	httpResp, err := o.Http.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode >= 400 {
		respBody, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("status code: %d, failed to read response body: %w", httpResp.StatusCode, err)
		}
		httpResp.Body.Close() // Ensure the body is closed
		return nil, errors.New(fmt.Sprintf("status code: %d, body: %s", httpResp.StatusCode, string(respBody)))
	}

	return httpResp, nil
}
