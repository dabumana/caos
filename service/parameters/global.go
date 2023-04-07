// Package parameters section
package parameters

import "caos/model"

// GlobalPreferences - General
type GlobalPreferences struct {
	// Engine properties
	Engine   string
	Mode     string
	Models   []string
	Roles    []string
	Template []string
	// Prompt properties
	Probabilities int32
	Results       int32
	MaxTokens     int64
	Temperature   float32
	Topp          float32
	Penalty       float32
	Frequency     float32
	PromptCtx     []string
	// Modes
	IsLoading         bool
	IsConversational  bool
	IsDeveloper       bool
	IsPredictable     bool
	IsEditable        bool
	IsTraining        bool
	IsNewSession      bool
	IsPromptReady     bool
	IsPromptStreaming bool
	IsTurbo           bool
	// Utilitaries
	Role       model.Roles
	CurrentID  string
	InlineText chan string
}

// ExternalBaseURL - External API endpoint
const ExternalBaseURL = "https://api.gptzero.me"
