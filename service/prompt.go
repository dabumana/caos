// Package service section
package service

import (
	"fmt"
	"log"

	"caos/service/parameters"

	"github.com/PullRequestInc/go-gpt3"
)

// Prompt - Handle prompt request
type Prompt struct {
	contextualResponse *gpt3.CompletionResponse
	extendedResponse   *gpt3.EditsResponse
}

// SendPrompt - Send task prompt
func (c Prompt) SendPrompt(service Client) *gpt3.CompletionResponse {
	if node.agent.currentUser.ctx == nil {
		log.Fatalln("Context NOT found")
	} else if service.client == nil {
		log.Fatalln("Client NOT found")
	}

	var prompt []string
	if parameters.IsConversational {
		prompt = []string{fmt.Sprintf("Human: %v \nAI:", service.promptProperties.PromptContext)}
	} else {
		prompt = service.promptProperties.PromptContext
	}

	req := gpt3.CompletionRequest{
		Prompt:           prompt,
		MaxTokens:        gpt3.IntPtr(service.promptProperties.MaxTokens),
		Temperature:      gpt3.Float32Ptr(service.engineProperties.Temperature),
		TopP:             gpt3.Float32Ptr(service.engineProperties.TopP),
		PresencePenalty:  *gpt3.Float32Ptr(service.engineProperties.PresencePenalty),
		FrequencyPenalty: *gpt3.Float32Ptr(service.engineProperties.FrequencyPenalty),
		Stream:           false,
		N:                gpt3.IntPtr(service.promptProperties.Results),
		LogProbs:         gpt3.IntPtr(service.promptProperties.Probabilities),
		Echo:             true}

	resp, err := service.client.CompletionWithEngine(
		node.agent.currentUser.ctx,
		service.engineProperties.Model,
		req)

	var event EventManager
	event.Errata(err)

	c.contextualResponse = resp
	return c.contextualResponse
}

// SendEditPrompt - Send edit instruction prompt
func (c Prompt) SendEditPrompt(service Client) *gpt3.EditsResponse {
	if service.ctx == nil {
		log.Fatalln("Context NOT found")
	} else if service.client == nil {
		log.Fatalln("Client NOT found")
	}

	req := gpt3.EditsRequest{
		Model:       service.engineProperties.Model,
		Input:       service.promptProperties.PromptContext[0],
		Instruction: service.promptProperties.Instruction[0],
		Temperature: gpt3.Float32Ptr(service.engineProperties.Temperature),
		TopP:        gpt3.Float32Ptr(service.engineProperties.TopP),
		N:           gpt3.IntPtr(service.promptProperties.Results)}

	resp, err := service.client.Edits(
		service.ctx,
		req)

	var event EventManager
	event.Errata(err)

	c.extendedResponse = resp
	return c.extendedResponse
}

// GetListModels - Get actual list of available models
func (c Prompt) GetListModels(service Client) *gpt3.EnginesResponse {
	resp, err := service.client.Engines(service.ctx)

	var event EventManager
	event.Errata(err)

	return resp
}
