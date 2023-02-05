// Package parameters section
package parameters

import "caos/model"

// PoolManager - Session pool management
type PoolManager struct {
	// EventPool - Historical events
	EventPool []model.HistoricalEvent
	// SessionPool - Historical sessions
	SessionPool []model.HistoricalSession
	// TrainingEventPool - Historical training events
	TrainingEventPool []model.HistoricalTrainingEvent
	// TrainingSessionPool - Historical training sessions
	TrainingSessionPool []model.HistoricalTrainingSession
}
