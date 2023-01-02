package model

// Engine preferences
type EngineProperties struct {
	Model            string
	Temperature      float32
	TopP             float32
	PresencePenalty  float32
	FrequencyPenalty float32
}
