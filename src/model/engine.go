// Package model section
package model

// Roles - Assignation roles
type Roles string

// const - Select Roles
const (
	System    Roles = "system"    // System role
	Assistant Roles = "assistant" // Assistant role
	User      Roles = "user"      // User role
)

// EngineProperties - Engine preferences
type EngineProperties struct {
	UserID           string  `json:"user_id"`
	Model            string  `json:"model"`
	Role             Roles   `json:"role"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"topp"`
	PresencePenalty  float32 `json:"presence_penalty"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
}
