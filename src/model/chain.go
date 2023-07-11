// Package model section
package model

// ChainPrompt - chain properties per prompt
type ChainPrompt struct {
	Source  []string `json:"source"`
	Context []string `json:"context"`
}

// ChainEvent - Session chain event
type ChainEvent struct {
	Timestamp string      `json:"timestamp"`
	Event     ChainPrompt `json:"event"`
}

// ChainSession - Chain session events
type ChainSession struct {
	ID      string       `json:"id"`
	Session []ChainEvent `json:"session"`
}
