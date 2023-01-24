// Package service section
package service

import (
	"fmt"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

// node - Global node service
var node Node

// Node - Node manager
type Node struct {
	prompt Prompt
	layout Layout
	agent  Controller
}

// ICommand - Command interface
type ICommand interface {
	AttachProfile() Client
	// StartInstructionRequest Controller API
	StartInstructionRequest() *gpt3.CompletionResponse
	StartRequest() *gpt3.CompletionResponse
}

// Start - Initialize node service
func (c Node) Start(sandboxMode bool) {
	var controller Controller
	var event EventManager
	node.agent = controller

	if node.agent.currentUser.client == nil {
		node.agent.currentUser = c.agent.AttachProfile()
	}

	if node.agent.currentUser.client != nil {
		event.LogClient(node.agent.currentUser)
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
