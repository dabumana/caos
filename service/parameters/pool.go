// Package parameters section
package parameters

import "caos/model"

// EventPool - Historical events
var EventPool []model.HistoricalEvent

// SessionPool - Historical sessions
var SessionPool []model.HistoricalSession

// TrainingEventPool - Historical training events
var TrainingEventPool []model.HistoricalTrainingEvent

// TrainingSessionPool - Historical training sessions
var TrainingSessionPool []model.HistoricalTrainingSession

// CurrentID - contextual parent id
var CurrentID string
