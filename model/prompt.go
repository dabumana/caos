// Prompt properties model
package model

// PromptProperties Prompt console preferences
type PromptProperties struct {
	PromptContext []string
	Prompt        []string
	MaxTokens     int
	Results       int
	Probabilities int
}
