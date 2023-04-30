// Package model section
package model

// TemplateProperties - Contextual template properties
type TemplateProperties struct {
	TemplateID string `json:"user_id"`
	// Prompt stages
	PromptContext        ChainPrompt `json:"prompt"`
	PromptAssemble       ChainPrompt `json:"assemble"`
	PromptImplementation ChainPrompt `json:"implementation"`
	PromptValidated      ChainPrompt `json:"validation"`
	// General properties per prompt
	Temperature float32 `json:"temperature"`
}
