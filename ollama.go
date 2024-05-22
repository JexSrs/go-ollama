package ollama

import (
	"bytes"
	json2 "encoding/json"
	"io"
	"net/http"
)

// Ollama represents a client for interacting with the Ollama API.
type Ollama struct {
	url          string
	defaultModel string
	chats        map[string]*Chat
	headers      map[string][]string

	Chat     ChatFunc
	Generate GenerateFunc

	BlobCheck  BlobCheckFunc
	BlobCreate BlobCreateFunc

	CreateModel     CreateModelFunc
	ListLocalModels ListLocalModelsFunc
	ShowModelInfo   ShowModelInfoFunc
	CopyModel       CopyModelFunc
	DeleteModel     DeleteModelFunc
	PullModel       PullModelFunc
	PushModel       PushModelFunc

	GenerateEmbeddings GenerateEmbeddingsFunc
}

// New creates a new Ollama client that points to the specified URL.
// It initializes the client with default settings and available API functions.
//
// Example:
//
//	llm := New("http://api.ollama.com")
func New(url string) *Ollama {
	o := &Ollama{
		url:          url,
		defaultModel: "llama3",
		headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		chats: make(map[string]*Chat),
	}

	o.Chat = o.newChatFunc()
	o.Generate = o.newGenerateFunc()

	o.BlobCheck = o.newBlobCheckFunc()
	o.BlobCreate = o.newBlobCreateFunc()

	o.CreateModel = o.newCreateModelFunc()
	o.ListLocalModels = o.newListLocalModelsFunc()
	o.ShowModelInfo = o.newShowModelInfoFunc()
	o.CopyModel = o.newCopyModelFunc()
	o.DeleteModel = o.newDeleteModelFunc()
	o.PullModel = o.newPullModelFunc()
	o.PushModel = o.newPushModelFunc()

	o.GenerateEmbeddings = o.newGenerateEmbeddingsFunc()

	return o
}

// Do makes an HTTP request to the specified path with the provided data.
// Pass a streamFunc for handling streaming types, or wait for the function to return the complete response.
//
// Parameters:
//   - path: The API endpoint path.
//   - data: The data to be sent in the request body, which will be marshaled to JSON.
//   - streamFunc: A function to handle streaming response chunks. If nil, the function waits for the complete response.
//
// Returns:
//   - A slice of byte slices containing the response data.
//   - An error if the request fails or if there is an issue reading the response.
//
// Example:
//
//	response, err := client.Do("/api/path", requestData, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("GCResponse:", response)
func (o *Ollama) Do(path string, data interface{}, streamFunc func(b []byte)) ([][]byte, error) {
	json, err := json2.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := o.Request("POST", path, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res [][]byte
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err == io.EOF {
			return nil, err
		}

		if n == 0 {
			break
		}

		res = append(res, buf[:n])

		if streamFunc != nil {
			streamFunc(buf[:n])
		}
	}

	return res, nil
}

// Request performs an HTTP request to the Ollama API.
//
// Parameters:
//   - method: The HTTP method (e.g., "GET", "POST").
//   - path: The API endpoint path.
//   - body: The request body as an io.Reader.
//
// Returns:
//   - An HTTP response from the API.
//   - An error if the request creation or execution fails.
func (o *Ollama) Request(method, path string, body io.Reader) (*http.Response, error) {
	httpReq, err := http.NewRequest(method, buildUrl(o.url, path), body)
	if err != nil {
		return nil, err
	}

	for k, v := range o.headers {
		httpReq.Header[k] = v
	}

	client := &http.Client{}
	return client.Do(httpReq)
}

// SetDefaultModel sets a default model to be used in requests if not specified.
// Defaults to "llama3".
//
// Parameters:
//   - v: The model name to set as the default.
func (o *Ollama) SetDefaultModel(v string) {
	o.defaultModel = v
}

// WithHeader sets additional headers to be included in requests.
//
// Parameters:
//   - key: The header key.
//   - value: The header values.
func (o *Ollama) WithHeader(key string, value []string) {
	o.headers[key] = value
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
