package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/gdamore/tcell/v2"
)

// Prompt service - handles requests prompts
type ServicePrompt struct {
	contextualResponse *gpt3.CompletionResponse
	extendedResponse   *gpt3.EditsResponse
}

/* Service prompt functionality */
// Log response details
func (c ServicePrompt) Log(resp *gpt3.CompletionResponse) {
	var responses []string
	for i := range resp.Choices {
		responses = append(responses, resp.Choices[i].Text)
	}
	promptctx = responses
	log := strings.Join(responses, "")
	reg := strings.ReplaceAll(log, "[]", "\n")
	Node.Layout.promptOutput.SetText(reg)
	Node.Layout.infoOutput.SetText(
		fmt.Sprintf("\nID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nToken probs: %v \nToken top: %v \n",
			resp.ID,
			resp.Model,
			resp.Created,
			resp.Object,
			resp.Usage.
				CompletionTokens,
			resp.Usage.PromptTokens,
			resp.Usage.TotalTokens,
			resp.Choices[0].FinishReason,
			resp.Choices[0].LogProbs.TokenLogprobs,
			resp.Choices[0].LogProbs.TopLogprobs))
}

// Log edited response details
func (c ServicePrompt) LogEdit(resp *gpt3.EditsResponse) {
	var responses []string
	for i := range resp.Choices {
		responses = append(responses, resp.Choices[i].Text, "\n")
	}
	promptctx = responses
	log := strings.Join(responses, "")
	reg := strings.ReplaceAll(log, "[]", "\n")
	Node.Layout.promptOutput.SetText(reg)
	Node.Layout.infoOutput.SetText(fmt.Sprintf("\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
		resp.Created,
		resp.Object,
		resp.Usage.CompletionTokens,
		resp.Usage.PromptTokens,
		resp.Usage.TotalTokens,
		resp.Choices[0].Index))
}

// Send taks prompt
func (c ServicePrompt) SendPrompt(service ServiceClient) *gpt3.CompletionResponse {
	if Node.Agent.currentUser.ctx == nil {
		log.Fatalln("Context NOT found")
	} else if service.client == nil {
		log.Fatalln("Client NOT found")
	}

	var prompt []string
	if isConversational {
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
		Stream:           true,
		N:                gpt3.IntPtr(service.promptProperties.Results),
		LogProbs:         gpt3.IntPtr(service.promptProperties.Probabilities),
		Echo:             true}

	resp, err := service.client.CompletionWithEngine(
		Node.Agent.currentUser.ctx,
		service.engineProperties.Model,
		req)

	if err != nil {
		Node.Layout.infoOutput.SetText(err.Error())
		Node.Layout.promptInput.SetPlaceholder("Press ENTER again to repeat the request.")
		Node.Layout.promptInput.SetPlaceholderTextColor(tcell.ColorDarkOrange)
	} else {
		Node.Layout.promptInput.SetPlaceholder("Type here...")
		Node.Layout.promptInput.SetPlaceholderTextColor(tcell.ColorBlack)
	}

	isLoading = false
	Node.Layout.promptInput.SetText("")

	c.contextualResponse = resp
	return c.contextualResponse
}

// Send instruction prompt
func (c ServicePrompt) SendIntructionPrompt(service ServiceClient) *gpt3.EditsResponse {
	if service.ctx == nil {
		log.Fatalln("Context NOT found")
	} else if service.client == nil {
		log.Fatalln("Client NOT found")
	}

	req := gpt3.EditsRequest{
		Model:       service.engineProperties.Model,
		Input:       service.promptProperties.PromptContext[0],
		Instruction: service.promptProperties.Prompt[0],
		Temperature: gpt3.Float32Ptr(service.engineProperties.Temperature),
		TopP:        gpt3.Float32Ptr(service.engineProperties.TopP),
		N:           gpt3.IntPtr(service.promptProperties.Results)}

	resp, err := service.client.Edits(
		service.ctx,
		req)

	if err != nil {
		Node.Layout.infoOutput.SetText(err.Error())
		Node.Layout.promptInput.SetPlaceholder("Press ENTER again to repeat the request.")
		Node.Layout.promptInput.SetPlaceholderTextColor(tcell.ColorDarkOrange)
	} else {
		Node.Layout.promptInput.SetPlaceholder("Type here...")
		Node.Layout.promptInput.SetPlaceholderTextColor(tcell.ColorBlack)
	}

	isLoading = false
	Node.Layout.promptInput.SetText("")

	c.extendedResponse = resp
	return c.extendedResponse
}
