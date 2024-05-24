package ollama

import (
	"time"
)

// ModelDetails represents detailed information about a model.
type ModelDetails struct {
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

// Message represents a message sent/received from the API.
type Message struct {
	Role    *string  `json:"role"`    // Role of the message, either system, user, or assistant.
	Content *string  `json:"content"` // Content of the message.
	Images  []string `json:"images"`  // Images associated with the message.
}

// Options represents the options that will be sent to the API.
type Options struct {
	NumKeep          *int     `json:"num_keep"`
	NumPredict       *int     `json:"num_predict"`       // Max number of tokens to predict.
	TopK             *int     `json:"top_k"`             // Reduces the probability of generating nonsense.
	TopP             *float64 `json:"top_p"`             // Controls diversity of text.
	TfsZ             *float64 `json:"tfs_z"`             // Tail free sampling.
	TypicalP         *float64 `json:"typical_p"`         // Typical probability.
	RepeatLastN      *int     `json:"repeat_last_n"`     // Prevents repetition.
	PenalizeNewLine  *bool    `json:"penalize_newline"`  // Penalizes new lines.
	RepeatPenalty    *float64 `json:"repeat_penalty"`    // Penalizes repetitions.
	PresencePenalty  *float64 `json:"presence_penalty"`  // Penalizes presence of tokens.
	FrequencyPenalty *float64 `json:"frequency_penalty"` // Penalizes frequency of tokens.
	Mirostat         *int     `json:"mirostat"`          // Enables Mirostat sampling.
	MirostatEta      *float64 `json:"mirostat_eta"`      // Learning rate for Mirostat.
	MirostatTau      *float64 `json:"mirostat_tau"`      // Balance between coherence and diversity.
	Stop             []string `json:"stop"`              // Stop sequences.
	Numa             *bool    `json:"numa"`              // NUMA support.
	NumCtx           *int     `json:"num_ctx"`           // Context window size.
	NumBatch         *int     `json:"num_batch"`         // Batch size.
	NumGPU           *int     `json:"num_gpu"`           // Number of GPUs.
	LowVRam          *bool    `json:"low_vram"`          // Low VRAM mode.
	F16KV            *bool    `json:"f16_kv"`            // 16-bit key-value pairs.
	VocabOnly        *bool    `json:"vocab_only"`        // Vocab only mode.
	NumThreads       *int     `json:"num_threads"`       // Number of threads.
	UseMMap          *bool    `json:"use_mmap"`          // Use memory-mapped files.
	UseMLock         *bool    `json:"use_mlock"`         // Use memory locking.
	Seed             *int     `json:"seed"`              // Random seed.
	Temperature      *float64 `json:"temperature"`       // Temperature for generation.
}

type Metrics struct {
	TotalDuration      time.Duration `json:"total_duration"`
	LoadDuration       time.Duration `json:"load_duration"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration"`
	EvalCount          int           `json:"eval_count"`
	EvalDuration       time.Duration `json:"eval_duration"`
}
