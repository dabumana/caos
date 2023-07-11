// Package parameters section
package parameters

import (
	"caos/model"
)

// GlobalPreferences - General
type GlobalPreferences struct {
	// Agent
	User     string
	Encoding string
	// Engine properties
	TemplateIDs       int
	Template          int
	Engine            string
	Mode              string
	Models            []string
	Roles             []string
	Historial         []model.HistoricalSession
	TrainingHistorial []model.TrainingSession
	// Prompt properties
	Probabilities int32
	Results       int32
	MaxTokens     int
	Temperature   float32
	Topp          float32
	Penalty       float32
	Frequency     float32
	PromptCtx     []string
	// Modes
	IsChained         bool
	IsLoading         bool
	IsEditable        bool
	IsNewSession      bool
	IsPromptReady     bool
	IsPromptStreaming bool
	// Utilitaries
	Role       model.Roles
	CurrentID  string
	InlineText chan string
}

// ExternalSearchBaseURL - External API endpoint
const ExternalSearchBaseURL = "https://www.google.com/search?client=firefox-b-m&gbv=1&q="
