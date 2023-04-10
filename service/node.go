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
	prompt     Prompt
	layout     Layout
	controller Controller
}

// Start - Initialize node service
func (c *Node) Start() {
	var controller Controller
	var event EventManager

	node.controller = controller

	if node.controller.currentAgent.client == nil {
		node.controller.currentAgent = c.controller.AttachProfile()
	}

	if node.controller.currentAgent.client != nil {
		event.LogClient(node.controller.currentAgent)
	} else {
		log.Fatalln("Client NOT loaded.")
		return
	}

	// Initialize app layout service
	InitializeLayout()

	// Exception
	if err := node.layout.app.Run(); err != nil {
		fmt.Printf("Execution error:%s\n", err)
	}
}
