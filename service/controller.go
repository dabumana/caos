// Service controller section
package service

import (
	"github.com/PullRequestInc/go-gpt3"
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
func (c Controller) InstructionRequest() *gpt3.EditsResponse {
	resp := node.prompt.SendInstructionPrompt(c.currentUser)

	var event EventManager
	c.currentUser.LogEngine()
	if resp != nil {
		event.LogEdit(&node.agent.currentUser.engineProperties, &node.agent.currentUser.promptProperties, resp)
		event.LogVizEdit(resp)
	}

	return resp
}

// CompletionRequest - Start completion request to send task prompt
func (c Controller) CompletionRequest() *gpt3.CompletionResponse {
	resp := node.prompt.SendPrompt(c.currentUser)

	var event EventManager
	c.currentUser.LogEngine()
	if resp != nil {
		event.Log(&node.agent.currentUser.engineProperties, &node.agent.currentUser.promptProperties, resp)
		event.LogViz(resp)
	}

	return resp
}
