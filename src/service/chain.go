// Package service section
package service

import (
	"caos/service/parameters"
	"fmt"

	"github.com/andelf/go-curl"
)

// chain - Sequence index
var seqChain []Chain

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
	return seqChain
}

// Set - Modify chain list
func Set(chain []Chain) {
	seqChain = chain
}

// Add - Add chain to actual stack
func (c *SequenceRequester) Add(chain *Chain) {
	seqChain = append(seqChain, *chain)
}

// Remove - Remove chain from actual stack
func (c *SequenceRequester) Remove(chain *Chain) {
	lastIndex := len(seqChain) - 1
	chain = &(seqChain)[lastIndex]
	seqChain = (seqChain)[:lastIndex]
}

// Chain - Event sequence object
type Chain struct {
	transform *Transformer
}

// Run - Run chain sequence
func (c *Chain) Run(service Agent, input []string) {
	c.transform.onConstructAssemble(service, input)
	c.transform.onConstructImplementation(service, input)
	c.transform.onConstructValidation(service, input)
}

// Transformer - LLM properties
type Transformer struct {
	model   string
	prompts []string
	stop    []string
}

// onConstructAssemble - Assemble transformer
func (c *Transformer) onConstructAssemble(service Agent, input []string) {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	req := fmt.Sprint(parameters.ExternalSearchBaseURL, input)

	easy.Setopt(curl.OPT_URL, req)
	easy.Setopt(curl.OPT_USERAGENT, service.preferences.User)
	easy.Setopt(curl.OPT_ACCEPT_ENCODING, service.preferences.Encode)
	easy.Setopt(curl.OPT_HTTP_CONTENT_DECODING, true)

	// Create a buffer for storing the response body
	var buf []byte
	writer := func(data []byte, userdata interface{}) bool {
		buf = append(buf, data...)
		return true
	}

	// Set write function to store response body in buffer
	easy.Setopt(curl.OPT_WRITEFUNCTION, writer)

	// Perform the request and check for errors
	_ = easy.Perform()

	fmt.Println(string(buf))
}

// onConstructImplementation - Implement transformer
func (c *Transformer) onConstructImplementation(service Agent, input []string) {}

// onConstructValidation - Valdiate transformer
func (c *Transformer) onConstructValidation(service Agent, input []string) {}
