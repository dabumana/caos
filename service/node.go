package service

import (
	"fmt"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

// Node - Global node service
var Node NodeService

// Global parameters
var engine string = "text-davinci-003"
var probabilities int32 = 1
var results int32 = 1
var temperature float32 = 0.4
var topp float32 = 1.0
var penalty float32 = 0.5
var frequency float32 = 0.5
var promptctx []string
var maxtokens int64 = 2048
var mode string = "Text"
var isLoading bool = false
var isConversational bool = false
var isEditable bool = false

// NodeService - Node manager
type NodeService struct {
	Prompt ServicePrompt
	Layout Layout
	Agent  Controller
}

// ICommand - Command interface
type ICommand interface {
	AttachProfile() Client
	// StartInstructionRequest Controller API
	StartInstructionRequest() *gpt3.CompletionResponse
	StartRequest() *gpt3.CompletionResponse
}

// Start - Initialize node service
func (c NodeService) Start(sandboxMode bool) {
	var controller Controller
	Node.Agent = controller

	if Node.Agent.currentUser.client == nil {
		Node.Agent.currentUser = c.Agent.AttachProfile()
	}

	if Node.Agent.currentUser.client != nil {
		Node.Agent.currentUser.LogClient()
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
	if err := Node.Layout.app.Run(); err != nil {
		fmt.Printf("Execution error:%s\n", err)
	}
}
