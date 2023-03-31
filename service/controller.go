// Package service section
package service

import "github.com/PullRequestInc/go-gpt3"

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
		event.LogEdit(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp)
		event.VisualLogEdit(resp)
	}

	event.LogEngine(node.controller.currentAgent)
}

// ChatCompletionRequest - Chat completion request to send task prompt
func (c Controller) ChatCompletionRequest() {
	var resp *gpt3.ChatCompletionResponse
	var cresp *gpt3.ChatCompletionStreamResponse
	if c.currentAgent.preferences.IsPromptStreaming {
		cresp = node.prompt.SendStreamingChatCompletion(c.currentAgent)
	} else {
		resp = node.prompt.SendChatCompletion(c.currentAgent)
	}

	var event EventManager
	if resp != nil && cresp == nil {
		event.LogChatCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp, nil)
		event.VisualLogChatCompletion(resp, nil)
	} else if cresp != nil && resp == nil {
		event.LogChatCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, nil, cresp)
		event.VisualLogChatCompletion(nil, cresp)
	}

	event.LogEngine(node.controller.currentAgent)
}

// CompletionRequest - Start completion request to send task prompt
func (c Controller) CompletionRequest() {
	var resp *gpt3.CompletionResponse
	if c.currentAgent.preferences.IsPromptStreaming {
		resp = node.prompt.SendStreamingCompletion(c.currentAgent)
	} else {
		resp = node.prompt.SendCompletion(c.currentAgent)
	}

	var event EventManager
	if resp != nil {
		event.LogCompletion(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp)
		event.VisualLogCompletion(resp)
	}

	event.LogEngine(node.controller.currentAgent)
}

// EmbeddingRequest - Start a embedding vector request
func (c Controller) EmbeddingRequest() {
	resp := node.prompt.SendEmbeddingPrompt(c.currentAgent)

	var event EventManager
	if resp != nil {
		event.LogEmbedding(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, resp)
		event.VisualLogEmbedding(resp)
	}

	event.LogEngine(node.controller.currentAgent)
}

// PredictableRequest - Start a predictable string request
func (c Controller) PredictableRequest() {
	resp := node.prompt.SendPredictablePrompt(c.currentAgent)

	var event EventManager
	if resp != nil {
		event.LogPredict(node.controller.currentAgent.engineProperties, node.controller.currentAgent.promptProperties, &node.controller.currentAgent.predictProperties, resp)
		event.VisualLogPredict(resp)
	}

	event.LogPredictEngine(node.controller.currentAgent)
}

// ListModels - Get actual models available
func (c Controller) ListModels() {
	resp := node.prompt.GetListModels(node.controller.currentAgent)
	if resp != nil {
		for _, i := range resp.Data {
			node.controller.currentAgent.preferences.Models = append(node.controller.currentAgent.preferences.Models, i.ID)
		}
	}
}
