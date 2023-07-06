// Package model section
package model

// TemplateProperties - Contextual template properties
type TemplateProperties struct {
	// Input
	Input []string `json:"input"`
	// Prompt stages
	PromptValidated ChainPrompt `json:"validation"`
}
