// Package service section
package service

// transform - Sequence index
var transform []Transformer

// TransformRequester - Transformer service interface
type TransformRequester []Transform

// TransformerSet - Transformer set component
type Transform interface {
	ClearListModel()
	GetListModel() []Transformer
	// Setup transformer
	SetupTransformer(llm string, input []string, stop []string) Transformer
}

func (c *TransformRequester) ClearListModel() {
	transform = nil
}

// GetListModel - Get actual list model
func (c *TransformRequester) GetListModel() []Transformer {
	return transform
}

// SetupTransformer - Add new element to listed interface
func (c *TransformRequester) SetupTransformer(llm string, input []string, stop []string) Transformer {
	chain := Transformer{
		model:   llm,
		prompts: input,
		stop:    stop,
	}

	transform = append(transform, chain)
	return chain
}

// Transformer - LLM properties
type Transformer struct {
	model   string
	prompts []string
	stop    []string
}

// onConstructAssemble - Assemble transformer
func (c *Transformer) onConstructAssemble() {}

// onConstructImplementation - Implement transformer
func (c *Transformer) onConstructImplementation() {}

// onConstructValidation - Valdiate transformer
func (c *Transformer) onConstructValidation() {}
