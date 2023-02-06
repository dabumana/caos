// Package model section
package model

// PoolProperties - Session pool management
type PoolProperties struct {
	Event           []HistoricalEvent           `json:"events"`
	Session         []HistoricalSession         `json:"sessions"`
	TrainingEvent   []HistoricalTrainingEvent   `json:"training_events"`
	TrainingSession []HistoricalTrainingSession `json:"training_sessions"`
}
