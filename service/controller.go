// Service controller section
package service

import "github.com/PullRequestInc/go-gpt3"

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
	c.currentUser.engineProperties.Model = "text-davinci-edit-001"
	resp := Node.Prompt.SendInstructionPrompt(c.currentUser)

	c.currentUser.LogEngine()
	if resp != nil {
		Node.Prompt.LogEdit(resp)
	}

	return resp
}

// CompletionRequest - Start completion request to send task prompt
func (c Controller) CompletionRequest() *gpt3.CompletionResponse {
	resp := Node.Prompt.SendPrompt(c.currentUser)

	c.currentUser.LogEngine()
	if resp != nil {
		Node.Prompt.Log(resp)
	}

	return resp
}
