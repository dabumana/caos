// Package service section
package service

import (
	"fmt"
	"log"
)

// node - Global node service
var node Node

// Node - Node manager
type Node struct {
	prompt Prompt
	layout Layout
	agent  Controller
}

// Start - Initialize node service
func (c Node) Start(sandboxMode bool) {
	var controller Controller
	var event EventManager
	node.agent = controller

	if node.agent.currentAgent.client == nil {
		node.agent.currentAgent = c.agent.AttachProfile()
	}

	if node.agent.currentAgent.client != nil {
		event.LogClient(node.agent.currentAgent)
	} else {
		log.Fatalln("Client NOT loaded.")
		return
	}
	// Initialize app layout service
	InitializeLayout()
	// Test mode
	if sandboxMode {
		return
	}
	// Exception
	if err := node.layout.app.Run(); err != nil {
		fmt.Printf("Execution error:%s\n", err)
	}
}
