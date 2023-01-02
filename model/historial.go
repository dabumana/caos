package model

type Historial struct {
	Id          string
	Created     int
	Object      string
	Requests    []Request
	Completions []Completion
}
