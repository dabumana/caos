// Engine properties model
package model

// EngineProperties - Engine preferences
type EngineProperties struct {
	Model            string  `json:"model"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"topp"`
	PresencePenalty  float32 `json:"presence_penalty"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
}
