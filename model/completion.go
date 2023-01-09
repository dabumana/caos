package model

type Completion struct {
	Id        string
	Created   int
	Object    string
	Response  []string
	TokenProb []string
}
