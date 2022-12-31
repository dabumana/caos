package service

import (
	"fmt"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

// Global node
var Node NodeService

// Global parameters
var engine string = "text-davinci-003"
var probabilities int = 1
var results int = 1
var temperature float32 = 0.0
var topp float32 = 1.0
var penalty float32 = 0.0
var frequency float32 = 0.0
var promptctx []string
var prompt []string
var maxtokens int = 150
var mode string = "Text"

// Node manager
type NodeService struct {
	Prompt ServicePrompt
	Layout ServiceLayout
	Agent  ServiceController
}

/* Command insterface */
type ICommand interface {
	AttachProfile() ServiceClient
	// Controller API
	StartInstructionRequest() *gpt3.CompletionResponse
	StartRequest() *gpt3.CompletionResponse
}

/** Node service functionality */
// Initialize node service
func (c NodeService) Start() {
	var controller ServiceController
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

	// Exception
	if err := Node.Layout.app.Run(); err != nil {
		fmt.Printf("Execution error:%s\n", err)
	}
}
