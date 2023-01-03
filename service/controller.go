package service

import "github.com/PullRequestInc/go-gpt3"

// Controller service - initializes the client as a service
type ServiceController struct {
	currentUser ServiceClient
}

/* Service controller functionality */
// Attach profile to a new service client
func (c ServiceController) AttachProfile() ServiceClient {
	var serviceClient ServiceClient
	serviceClient = serviceClient.Initialize()

	return serviceClient
}

const TextDavinci001Edit string = "text-davinci-edit-001"

// Start edited request based on the previous response
func (c ServiceController) InstructionRequest() *gpt3.EditsResponse {
	var servicePrompt ServicePrompt
	c.currentUser.engineProperties.Model = TextDavinci001Edit

	resp := servicePrompt.SendIntructionPrompt(c.currentUser)

	c.currentUser.LogEngine()
	if resp != nil {
		Node.Prompt.LogEdit(resp)
	}

	return resp
}

// Start initial request to send task prompt
func (c ServiceController) StartRequest() *gpt3.CompletionResponse {
	resp := Node.Prompt.SendPrompt(c.currentUser)

	c.currentUser.LogEngine()
	if resp != nil {
		Node.Prompt.Log(resp)
	}

	return resp
}
