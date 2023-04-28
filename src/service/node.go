// Package service section
package service

import (
	"log"
)

// node - Global node service
var node *Node

// Node - Node manager
type Node struct {
	prompt     Prompt
	layout     Layout
	controller Controller
}

// Init - Entrypoint for terminal service
func (c *Node) Init() {
	node = onConstruct()
	node.Start()
}

// onConstruct - Construct agent
func onConstruct() *Node {
	service := Node{}
	return &service
}

// Start - Initialize node service
func (c *Node) Start() {
	var controller Controller
	var event EventManager

	c.controller = controller

	if c.controller.currentAgent.client == nil {
		c.controller.currentAgent = c.controller.AttachProfile()
	}

	if c.controller.currentAgent.client != nil {
		event.LogClient(c.controller.currentAgent)
	} else {
		log.Fatalln("Client NOT loaded.")
		return
	}
	// Generate service
	c.layout.app, c.layout.screen = ConstructService()
	// Initialize app layout service
	InitializeService()

}
