// Package service section
package service

import (
	"caos/service/parameters"
)

// Controller - Contextual client controller API
type Controller struct {
	currentUser Client
}

// AttachProfile - Attach profile to a new service client
func (c Controller) AttachProfile() Client {
	var serviceClient Client
	serviceClient = serviceClient.Initialize()

	return serviceClient
}

// InstructionRequest - Start edit request  to send a task prompt
func (c Controller) InstructionRequest() {
	resp := node.prompt.SendInstructionPrompt(c.currentUser)

	var event EventManager
	if resp != nil {
		event.LogInstruction(&node.agent.currentUser.engineProperties, &node.agent.currentUser.promptProperties, resp)
		event.VisualLogInstruction(resp)
	}
	event.LogEngine(c.currentUser)
}

// CompletionRequest - Start completion request to send task prompt
func (c Controller) CompletionRequest() {
	resp := node.prompt.SendPrompt(c.currentUser)

	var event EventManager
	if resp != nil {
		event.LogCompletion(&node.agent.currentUser.engineProperties, &node.agent.currentUser.promptProperties, resp)
		event.VisualLogCompletion(resp)
	}
	event.LogEngine(c.currentUser)
}

// ListModels - Get actual models available
func (c Controller) ListModels() {
	resp := node.prompt.GetListModels(c.currentUser)
	for _, i := range resp.Data {
		parameters.Models = append(parameters.Models, i.ID)
	}
}
