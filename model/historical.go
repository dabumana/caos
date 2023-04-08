// Package model section
package model

// HistoricalPrompt - historical prompt
type HistoricalPrompt struct {
	Header         EngineProperties  `json:"properties"`
	Body           PromptProperties  `json:"parameters"`
	PredictiveBody PredictProperties `json:"predictivity"`
}

// HistoricalEvent - Session historical event
type HistoricalEvent struct {
	Timestamp string           `json:"timestamp"`
	Event     HistoricalPrompt `json:"event"`
}

// HistoricalSession - Historical session events
type HistoricalSession struct {
	ID      string            `json:"id"`
	Session []HistoricalEvent `json:"session"`
}
