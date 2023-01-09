package model

type Historical struct {
	Id          string
	Created     int
	Object      string
	Requests    []Request
	Completions []Completion
}
