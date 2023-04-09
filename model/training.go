// Package model section
package model

// TrainingPrompt - Fine-tunning model
type TrainingPrompt struct {
	Prompt     []string `json:"prompt"`
	Completion []string `json:"completion"`
}

// TrainingEvent - Session training historical event
type TrainingEvent struct {
	Timestamp string         `json:"timestamp"`
	Event     TrainingPrompt `json:"event"`
}

// TrainingSession - Historical training session events
type TrainingSession struct {
	ID      string          `json:"id"`
	Session []TrainingEvent `json:"session"`
}
