// Package model section
package model

// Roles - Assignation roles
type Roles string

const (
	System    Roles = "system"
	Assistant Roles = "assistant"
	User      Roles = "user"
)

// EngineProperties - Engine preferences
type EngineProperties struct {
	UserId           string  `json:"user_id"`
	Model            string  `json:"model"`
	Role             Roles   `json:"role"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"topp"`
	PresencePenalty  float32 `json:"presence_penalty"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
}
