// Interaction prompt models
package model

// TrainingPrompt - Fine-tunning model
type TrainingPrompt struct {
	Prompt     []string `json:"prompt"`
	Completion []string `json:"completion"`
}
