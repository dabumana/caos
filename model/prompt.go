// Prompt properties model
package model

// PromptProperties - Prompt console preferences
type PromptProperties struct {
	PromptContext []string `json:"prompt"`
	Instruction   []string `json:"instruction"`
	MaxTokens     int      `json:"token_ammount"`
	Results       int      `json:"results"`
	Probabilities int      `json:"probabilities"`
}
