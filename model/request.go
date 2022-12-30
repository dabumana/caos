package model

type Request struct {
	Id                string
	Created           int
	Object            string
	Prompt            []string
	Instruction       []string
	Result            string
	EnginePreferences []EngineProperties
}
