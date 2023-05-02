// Package service section
package service

// ChainedTransforms - Sequence index
var chainedTransform []Chain

// SequenceRequester - Chained Transform template
type SequenceRequester []Sequence

// Sequence - Chain sequence interface
type Sequence interface {
	Get() []Chain
	Set(chainedTransform []Chain)
	// Stack operation
	Add(chain *Chain)
	Remove(chain *Chain)
}

// Get - Retrieve actual chain list
func (c *SequenceRequester) Get() []Chain {
	return chainedTransform
}

// Set - Modify chain list
func Set(chain []Chain) {
	chainedTransform = chain
}

// Add - Add chain to actual stack
func (c *SequenceRequester) Add(chain *Chain) {
	chainedTransform = append(chainedTransform, *chain)
}

// Remove - Remove chain from actual stack
func (c *SequenceRequester) Remove(chain *Chain) {
	lastIndex := len(chainedTransform) - 1
	chain = &(chainedTransform)[lastIndex]
	chainedTransform = (chainedTransform)[:lastIndex]
}

// Chain - Event sequence object
type Chain struct {
	transforms []Transformer
}

// PrepareChain - Initialize chain sequence
func (c *Chain) PrepareChain() {
	for i := range c.transforms {
		c.transforms[i].onConstructAssemble()
		c.transforms[i].onConstructImplementation()
		c.transforms[i].onConstructValidation()
	}
}

// RunSequence - Run chain sequence
func (c *Chain) RunSequence(input string) {}
