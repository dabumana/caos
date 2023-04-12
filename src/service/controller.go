// Package service section
package service

import (
	"caos/model"

	"github.com/PullRequestInc/go-gpt3"
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

// EditRequest - Start edit request  to send a task prompt
func (c *Controller) EditRequest() {
	resp := node.prompt.SendEditPrompt(c.currentAgent)

	if resp != nil {
		c.events.LogEdit(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp)
		c.events.VisualLogEdit(resp)
	}

	c.events.LogEngine(c.currentAgent)
}

// ChatCompletionRequest - Chat completion request to send task prompt
func (c *Controller) ChatCompletionRequest() {
	var resp *gpt3.ChatCompletionResponse
	var cresp *gpt3.ChatCompletionStreamResponse
	if c.currentAgent.preferences.IsPromptStreaming {
		cresp = node.prompt.SendStreamingChatCompletion(c.currentAgent)
	} else {
		resp = node.prompt.SendChatCompletion(c.currentAgent)
	}

	if resp != nil && cresp == nil {
		c.events.LogChatCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp, nil)
		c.events.VisualLogChatCompletion(resp, nil)
	} else if cresp != nil && resp == nil {
		c.events.LogChatCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, nil, cresp)
		c.events.VisualLogChatCompletion(nil, cresp)
	}

	c.events.LogEngine(c.currentAgent)
}

// CompletionRequest - Start completion request to send task prompt
func (c *Controller) CompletionRequest() {
	var resp *gpt3.CompletionResponse
	if c.currentAgent.preferences.IsPromptStreaming {
		resp = node.prompt.SendStreamingCompletion(c.currentAgent)
	} else {
		resp = node.prompt.SendCompletion(c.currentAgent)
	}

	if resp != nil {
		c.events.LogCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp)
		c.events.VisualLogCompletion(resp)
	}

	c.events.LogEngine(c.currentAgent)
}

// EmbeddingRequest - Start a embedding vector request
func (c *Controller) EmbeddingRequest() {
	resp := node.prompt.SendEmbeddingPrompt(c.currentAgent)

	if resp != nil {
		c.events.LogEmbedding(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp)
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
