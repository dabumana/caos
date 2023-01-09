package model

type Request struct {
	Id                string
	Created           int
	PromptPreferences []PromptProperties
	EnginePreferences []EngineProperties
}
