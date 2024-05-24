package ollama

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

var LLM *Ollama

const (
	_Model   = "phi3"
	_Message = "Respond only with the text \"Hello I am here to assist you.\""
)

func init() {
	uri, _ := url.Parse("http://localhost:11434")
	LLM = New(*uri)
}

func TestGenerateStream(t *testing.T) {
	streamedResponses := make([]GenerateResponse, 0)
	resp, err := LLM.Generate(
		LLM.Generate.WithModel(_Model),
		LLM.Generate.WithSeed(123),
		LLM.Generate.WithPrompt(_Message),
		LLM.Generate.WithKeepAlive("5m"),
		LLM.Generate.WithTemperature(0.1),
		LLM.Generate.WithStream(true, 512000, func(r *GenerateResponse, err error) {
			if err != nil {
				t.Errorf("Generate streaming returned an error: %s", err)
			} else {
				streamedResponses = append(streamedResponses, *r)
			}
		}),
	)

	if err != nil {
		t.Errorf("Generate returned an error: %s", err)
		return
	}

	result := ""
	for _, response := range streamedResponses {
		result += response.Response
	}

	if resp.Response != result {
		t.Errorf("Expected \"%s\", got \"%s\"", resp.Response, result)
		return
	}

	resp, err = LLM.Generate(
		LLM.Generate.WithModel(_Model),
		LLM.Generate.WithSeed(123),
		LLM.Generate.WithPrompt(_Message),
		LLM.Generate.WithKeepAlive("5m"),
		LLM.Generate.WithTemperature(0.1),
	)

	if err != nil {
		t.Errorf("Generate returned an error: %s", err)
		return
	}

	if resp.Response != result {
		t.Errorf("Expected \"%s\", got \"%s\"", resp.Response, result)
		return
	}
}

func TestChatStream(t *testing.T) {
	streamedResponses := make([]ChatResponse, 0)
	resp, err := LLM.Chat(
		nil,
		LLM.Chat.WithModel(_Model),
		LLM.Chat.WithSeed(123),
		LLM.Chat.WithMessage(Message{Role: pointer("user"), Content: pointer(_Message)}),
		LLM.Chat.WithKeepAlive("5m"),
		LLM.Chat.WithTemperature(0.1),
		LLM.Chat.WithStream(true, 512000, func(r *ChatResponse, err error) {
			if err != nil {
				t.Errorf("Chat streaming returned an error: %s", err)
			} else {
				streamedResponses = append(streamedResponses, *r)
			}
		}),
	)

	if err != nil {
		t.Errorf("Chat returned an error: %s", err)
		return
	}

	result := ""
	for _, response := range streamedResponses {
		result += *response.Message.Content
	}

	if *resp.Message.Content != result {
		t.Errorf("Expected \"%s\", got \"%s\"", *resp.Message.Content, result)
		return
	}

	resp, err = LLM.Chat(
		nil,
		LLM.Chat.WithModel(_Model),
		LLM.Chat.WithSeed(123),
		LLM.Chat.WithMessage(Message{Role: pointer("user"), Content: pointer(_Message)}),
		LLM.Chat.WithKeepAlive("5m"),
		LLM.Chat.WithTemperature(0.1),
	)

	if err != nil {
		t.Errorf("Chat returned an error: %s", err)
		return
	}

	if *resp.Message.Content != result {
		t.Errorf("Expected \"%s\", got \"%s\"", *resp.Message.Content, result)
		return
	}
}

func TestBlobs(t *testing.T) {
	data := []byte{0x0, 0x1, 0xF}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(data))

	err := LLM.Blobs.Create(digest, data)
	if err != nil {
		t.Errorf("Blobs returned an error: %s", err)
		return
	}

	err = LLM.Blobs.Check(digest)
	if err != nil {
		t.Errorf("Blobs returned an error: %s", err)
		return
	}
}

func TestModelsActions(t *testing.T) {
	err := LLM.Models.Copy(_Model, _Model+"-test")
	if err != nil {
		t.Errorf("Models returned an error: %s", err)
		return
	}

	models, err := LLM.Models.List()
	if err != nil {
		t.Errorf("Models returned an error: %s", err)
		return
	}

	found := false
	for _, model := range models.Models {
		if strings.HasPrefix(model.Name, _Model+"-test") {
			found = true
		}
	}

	if !found {
		t.Errorf("Models did not find test model")
		return
	}

	_, err = LLM.Models.ShowInfo(
		LLM.Models.ShowInfo.WithModel(_Model + "-test"),
	)

	if err != nil {
		t.Errorf("Models returned an error: %s", err)
		return
	}

	err = LLM.Models.Delete(_Model + "-test")
	if err != nil {
		t.Errorf("Models returned an error: %s", err)
		return
	}
}
