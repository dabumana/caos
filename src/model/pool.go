// Package model section
package model

// PoolProperties - Session pool management
type PoolProperties struct {
	Event           []HistoricalEvent   `json:"events"`
	Session         []HistoricalSession `json:"sessions"`
	TrainingEvent   []TrainingEvent     `json:"training_events"`
	TrainingSession []TrainingSession   `json:"training_sessions"`
}
