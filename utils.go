package ollama

import (
	json2 "encoding/json"
	"strings"
)

func bodyTo[T any](body []byte) (*T, error) {
	var response T
	err := json2.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func pointer[T any](t T) *T {
	return &t
}

func buildUrl(baseUrl, path string) string {
	url := baseUrl

	if !strings.HasSuffix(url, "/") && !strings.HasPrefix(path, "/") {
		url += "/"
	} else if strings.HasSuffix(url, "/") && strings.HasPrefix(path, "/") {
		url = url[0 : len(url)-1]
	}

	url += path
	return url
}
