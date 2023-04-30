// Package service section
package service

import (
	"github.com/tmc/langchaingo/llms"
)

// Transformer - LLM properties
type Transformer struct {
	model   *llms.LLM
	prompts []string
	stop    []string
}

// SetupTransformer
func (c *Transformer) SetupTransformer(llm *llms.LLM, input []string, stop []string) {
	c.model = llm
	c.prompts = input
	c.stop = stop
}
