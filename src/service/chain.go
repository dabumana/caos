// Package service section
package service

// Chain - Event sequence object
type Chain struct {
	transforms []Transformer
}

// onConstructAssemble - Assemble transformer
func (c *Transformer) onConstructAssemble() {}

// onConstructImplementation - Implement transformer
func (c *Transformer) onConstructImplementation() {}

// onConstructValidation - Valdiate transformer
func (c *Transformer) onConstructValidation() {}

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
