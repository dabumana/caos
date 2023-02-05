// Package parameters section
package parameters

// GlobalPreferences - General
type GlobalPreferences struct {
	// Engine properties
	Engine string
	Mode   string
	Models []string
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
	IsLoading        bool
	IsConversational bool
	IsPredictable    bool
	IsEditable       bool
	IsTraining       bool
	IsNewSession     bool
	IsPromptReady    bool
	// Utilitaries
	CurrentID string
}

const ExternalBaseURL = "https://api.gptzero.me"
