// Package service section
package service

import (
	"caos/service/parameters"
)

// Controller - Contextual client controller API
type Controller struct {
	currentAgent Agent
}

// AttachProfile - Attach profile to a new service client
func (c Controller) AttachProfile() Agent {
	var serviceClient Agent
	serviceClient = serviceClient.Initialize()

	return serviceClient
}

// EditRequest - Start edit request  to send a task prompt
func (c Controller) EditRequest() {
	resp := node.prompt.SendEditPrompt(c.currentAgent)

	var event EventManager
	if resp != nil {
		event.LogEdit(&node.agent.currentAgent.engineProperties, &node.agent.currentAgent.promptProperties, resp)
		event.VisualLogEdit(resp)
	}
	event.LogEngine(c.currentAgent)
}

// CompletionRequest - Start completion request to send task prompt
func (c Controller) CompletionRequest() {
	resp := node.prompt.SendCompletion(c.currentAgent)

	var event EventManager
	if resp != nil {
		event.LogCompletion(&node.agent.currentAgent.engineProperties, &node.agent.currentAgent.promptProperties, resp)
		event.VisualLogCompletion(resp)
	}
	event.LogEngine(c.currentAgent)
}

// EmbeddingRequest - Start a embedding vector request
func (c Controller) EmbeddingRequest() {
	resp := node.prompt.SendEmbeddingPrompt(c.currentAgent)

	var event EventManager
	if resp != nil {
		event.LogEmbedding(&node.agent.currentAgent.engineProperties, &node.agent.currentAgent.promptProperties, resp)
		event.VisualLogEmbedding(resp)
	}
	event.LogEngine(c.currentAgent)
}

// PredictableRequest - Start a predictable string request
func (c Controller) PredictableRequest() {
	resp := node.prompt.SendPredictablePrompt(c.currentAgent)

	var event EventManager
	if resp != nil {
		event.LogPredict(&node.agent.currentAgent.predictProperties, resp)
		event.VisualLogPredict(resp)
	}
	event.LogPredictEngine(node.agent.currentAgent)
}

// ListModels - Get actual models available
func (c Controller) ListModels() {
	resp := node.prompt.GetListModels(c.currentAgent)
	for _, i := range resp.Data {
		parameters.Models = append(parameters.Models, i.ID)
	}
}
