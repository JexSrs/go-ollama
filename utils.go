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

func splitJSONObjects(data []byte) [][]byte {
	var results [][]byte
	var stack []byte
	var start, end int
	var inString bool

	for i := 0; i < len(data); i++ {
		switch data[i] {
		case '{':
			if !inString {
				if len(stack) == 0 {
					start = i
				}
				stack = append(stack, '{')
			}
		case '}':
			if !inString {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 {
					end = i + 1
					results = append(results, data[start:end])
				}
			}
		case '"':
			if i == 0 || data[i-1] != '\\' {
				inString = !inString
			}
		}
	}

	return results
}
