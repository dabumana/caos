// Package model section
package model

// ChainPrompt - chain properties per prompt
type ChainPrompt struct {
	Model      string         `json:"model"`
	Text       []string       `json:"input"`
	Generative map[string]any `json:"generative"`
}

// ChainEvent - Session chain event
type ChainEvent struct {
	Timestamp string      `json:"timestamp"`
	Event     ChainPrompt `json:"event"`
}

// ChainSession - Chain session events
type ChainSession struct {
	ID      string            `json:"id"`
	Session []HistoricalEvent `json:"session"`
}
