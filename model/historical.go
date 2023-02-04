// Package model section
package model

// HistoricalPrompt - historical prompt
type HistoricalPrompt struct {
	Header  EngineProperties  `json:"properties"`
	Body    PromptProperties  `json:"body"`
	Predict PredictProperties `json:"predict"`
}

// HistoricalEvent - Session historical event
type HistoricalEvent struct {
	Timestamp string           `json:"timestamp"`
	Event     HistoricalPrompt `json:"event"`
}

// HistoricalTrainingEvent - Session training historical event
type HistoricalTrainingEvent struct {
	Timestamp string         `json:"timestamp"`
	Event     TrainingPrompt `json:"event"`
}

// HistoricalSession - Historical session events
type HistoricalSession struct {
	ID      string            `json:"id"`
	Session []HistoricalEvent `json:"session"`
}

// HistoricalTrainingSession - Historical training session events
type HistoricalTrainingSession struct {
	ID      string                    `json:"id"`
	Session []HistoricalTrainingEvent `json:"session"`
}
