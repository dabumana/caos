// Package model section
package model

// TemplateProperties - Contextual template properties
type TemplateProperties struct {
	TemplateID int `json:"user_id"`
	// Input
	Input []string `json:"input"`
	// Prompt stages
	PromptAssemble       ChainPrompt `json:"assemble"`
	PromptImplementation ChainPrompt `json:"implementation"`
	PromptValidated      ChainPrompt `json:"validation"`
	// General properties per prompt
	Temperature float32 `json:"temperature"`
}
