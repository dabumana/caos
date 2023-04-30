// Package service section
package service

// ChainedTransforms - Sequence index
var chainedTransform []Chain

// SequenceHandler - Chained Transform template
type SequenceHandler []ChainTransforms

// SequencesInterface - Chain sequence interface
type ChainTransforms interface {
	Get() []Chain
	Set(chainedTransform []Chain)
	// Stack operation
	Add(chain *Chain)
	Remove(chain *Chain)
}

// Get - Retrieve actual chain list
func (c *SequenceHandler) Get() []Chain {
	return chainedTransform
}

// Set - Modify chain list
func Set(chain []Chain) {
	chainedTransform = chain
}

// Add - Add chain to actual stack
func (c *SequenceHandler) Add(chain *Chain) {
	chainedTransform = append(chainedTransform, *chain)
}

// Remove - Remove chain from actual stack
func (c *SequenceHandler) Remove(chain *Chain) {
	lastIndex := len(chainedTransform) - 1
	chain = &(chainedTransform)[lastIndex]
	chainedTransform = (chainedTransform)[:lastIndex]
}

// onConstructAssemble - Assemble transformer
func (c *Transformer) onConstructAssemble() {}

// onConstructImplementation - Implement transformer
func (c *Transformer) onConstructImplementation() {}

// onConstructValidation - Valdiate transformer
func (c *Transformer) onConstructValidation() {}

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
func (c *Chain) RunSequence() {}
