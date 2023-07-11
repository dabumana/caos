// Package model section
package model

// PromptProperties - Prompt console preferences
type PromptProperties struct {
	Input         []string `json:"prompt"`
	Instruction   []string `json:"instruction"`
	Content       []string `json:"content"`
	MaxTokens     int      `json:"token_ammount"`
	Results       int      `json:"results"`
	Probabilities int      `json:"probabilities"`
}
