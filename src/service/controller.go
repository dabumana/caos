// Package service section
package service

import (
	"caos/model"
)

// Controller - Contextual client controller API
type Controller struct {
	currentAgent Agent
	events       EventManager
}

// AttachProfile - Attach profile to a new service client
func (c *Controller) AttachProfile() Agent {
	var serviceClient Agent
	serviceClient = serviceClient.Initialize()

	return serviceClient
}

// FlushEvents - Reset the pool
func (c *Controller) FlushEvents() {
	node.controller.events.pool.TrainingEvent = []model.TrainingEvent{}
	node.controller.events.pool.TrainingSession = []model.TrainingSession{}
}

// ChatCompletionRequest - Chat completion request to send task prompt
func (c *Controller) ChatCompletionRequest() {
	if c.currentAgent.preferences.IsPromptStreaming {
		resp, _ := node.prompt.SendChatCompletion(c.currentAgent)

		c.events.LogChatCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, nil, resp)
		c.events.VisualLogCompletion(nil, nil, resp)
	} else {
		_, resp := node.prompt.SendChatCompletion(c.currentAgent)

		if resp != nil {
			c.events.LogChatCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp, nil)
			c.events.VisualLogCompletion(nil, resp, nil)
		}
	}

	c.events.LogEngine(c.currentAgent)
}

// CompletionRequest - Start completion request to send task prompt
func (c *Controller) CompletionRequest() {
	resp := node.prompt.SendCompletion(c.currentAgent)

	if resp != nil {
		c.events.LogGeneralCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, []string{resp.Choices[0].Text}, resp.ID)
		c.events.VisualLogCompletion(resp, nil, nil)
	}

	c.events.LogEngine(c.currentAgent)
}

// EditRequest - Start edit request  to send a task prompt
func (c *Controller) EditRequest() {
	resp := node.prompt.SendEditPrompt(c.currentAgent)

	if resp != nil {
		c.events.LogGeneralCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, []string{resp.Choices[0].Text}, node.controller.currentAgent.preferences.CurrentID)
		c.events.VisualLogEdit(resp)
	}

	c.events.LogEngine(c.currentAgent)
}

// EmbeddingRequest - Start a embedding vector request
func (c *Controller) EmbeddingRequest() {
	resp := node.prompt.SendEmbeddingPrompt(c.currentAgent)

	if resp != nil {
		c.events.LogGeneralCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, []string{resp.Data[0].Object}, node.controller.currentAgent.preferences.CurrentID)
		c.events.VisualLogEmbedding(resp)
	}

	c.events.LogEngine(c.currentAgent)
}

// PredictableRequest - Start a predictable string request
func (c *Controller) PredictableRequest() {
	resp := node.prompt.SendPredictablePrompt(c.currentAgent)

	if resp != nil {
		c.events.LogPredict(node.controller.currentAgent.engineProperties, node.controller.currentAgent.predictProperties, resp)
		c.events.VisualLogPredict(resp)

		c.currentAgent.predictProperties.Details.Documents = append(c.currentAgent.predictProperties.Details.Documents, resp.Documents...)
	}

	c.events.LogPredictEngine(c.currentAgent)
}

// ListModels - Get actual models available
func (c *Controller) ListModels() {
	resp := node.prompt.GetListModels(c.currentAgent)
	if resp != nil {
		for _, i := range resp.Data {
			node.controller.currentAgent.preferences.Models = append(node.controller.currentAgent.preferences.Models, i.ID)
		}
	}
}
