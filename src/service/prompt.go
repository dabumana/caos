// Package service section
package service

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"caos/model"
	"caos/util"
	"encoding/json"
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/mitchellh/go-wordwrap"
	"github.com/rivo/tview"
)

// Prompt - Handle prompt request
type Prompt struct {
	contextualResponse *gpt3.CompletionResponse
	extendedResponse   *gpt3.EditsResponse
	embeddingResponse  *gpt3.EmbeddingsResponse
	// Chat completion response
	chatStreamResponse  *gpt3.ChatCompletionStreamResponse
	chatResponse        *gpt3.ChatCompletionResponse
	predictableResponse *model.PredictResponse
}

// isContextValid - Client context validation
func isContextValid(current Agent) bool {
	if current.ctx == nil {
		log.Fatalln("Context NOT found")
		return false
	} else if current.client == nil {
		log.Fatalln("Client NOT found")
		return false
	}

	return true
}

// isTestingEnvironment - Validate when test executes
func isTestingEnvironment() bool {
	hasTestFlag := flag.Lookup("test.v") != nil
	return hasTestFlag
}

// SendChatCompletionPrompt - Send streaming chat completion prompt
func (c *Prompt) SendChatCompletionPrompt(service Agent) (*gpt3.ChatCompletionStreamResponse, *gpt3.ChatCompletionResponse) {
	if isContextValid(service) {
		var buffer []string

		prompt := service.SetTemplate(string(""), service.PromptProperties.Input[0])[0]
		urls, ctxVerified := service.SetContext(&service.PromptProperties)

		service.TemplateProperties.PromptValidated.Source = append(service.TemplateProperties.PromptValidated.Source, urls...)
		service.TemplateProperties.PromptValidated.Context = append(service.TemplateProperties.PromptValidated.Context, ctxVerified...)

		msg := fmt.Sprint(
			prompt,
			"\nNow you have real-time access to the internet, I'll provide some contextual information with urls for main reference for your responses",
			"\nthis is fundamentally from real-time online actual results, provide a response around the following contextual information:",
			ctxVerified,
			"\nObtained from the following urls:",
			urls,
			"\nPlease always elaborate a detailed response with the following complete schema (KEEP LINE BY LINE):",
			"\nQuestion: <User input ONLY>",
			"\nResponse: <Your Detailed response should contain the contextual information from the real-time online results ONLY>",
			"\nResume: <Include more than 3000-WORDS per response ONLY>",
			"\nSuggestions: <based entirely on the verified context include suggestions to look or search ONLY>",
			"\nSource: <List all the urls from the contextual information ONLY>",
		)

		if !isTestingEnvironment() {
			service.PromptProperties.MaxTokens = len(util.EncodePromptBytePair([]string{msg}, service.EngineProperties.Model))
			node.controller.currentAgent.preferences.MaxTokens = service.PromptProperties.MaxTokens
		}

		ctx := gpt3.ChatCompletionRequestMessage{
			Role:    string(service.preferences.Role),
			Content: msg,
		}

		req := gpt3.ChatCompletionRequest{
			Model:            service.EngineProperties.Model,
			User:             service.id,
			Messages:         []gpt3.ChatCompletionRequestMessage{ctx},
			MaxTokens:        *gpt3.IntPtr(service.PromptProperties.MaxTokens),
			Temperature:      *gpt3.Float32Ptr(service.EngineProperties.Temperature),
			TopP:             *gpt3.Float32Ptr(service.EngineProperties.TopP),
			PresencePenalty:  *gpt3.Float32Ptr(service.EngineProperties.PresencePenalty),
			FrequencyPenalty: *gpt3.Float32Ptr(service.EngineProperties.FrequencyPenalty),
			Stream:           service.preferences.IsPromptStreaming,
			N:                *gpt3.IntPtr(service.PromptProperties.Results),
			Stop:             []string{"stop"},
		}

		if service.preferences.IsPromptStreaming {
			sresp := &gpt3.ChatCompletionStreamResponse{}

			var bWriter tview.TextViewWriter
			if !isTestingEnvironment() {
				bWriter = node.layout.promptOutput.BatchWriter()
				defer bWriter.Close()
				bWriter.Clear()
				bWriter.Write([]byte("\n"))
				buffer = append(buffer, "\n")
			}

			fmt.Print("\033[H\033[2J")
			client := *service.client
			err := client.ChatCompletionStream(
				service.ctx,
				req, func(out *gpt3.ChatCompletionStreamResponse) {
					sresp.ID = out.ID
					sresp.Choices = out.Choices
					sresp.Created = out.Created
					sresp.Model = out.Model
					sresp.Object = out.Object
					sresp.Usage = out.Usage
					// Choices
					sresp.Choices[0].Index = out.Choices[0].Index
					sresp.Choices[0].FinishReason = out.Choices[0].FinishReason
					// Delta
					sresp.Choices[0].Delta = out.Choices[0].Delta
					sresp.Choices[0].Delta.Content = out.Choices[0].Delta.Content
					sresp.Choices[0].Delta.Role = out.Choices[0].Delta.Role
					// Write buffer
					if !isTestingEnvironment() {
						buffer = append(buffer, out.Choices[0].Delta.Content)
						bWriter.Write([]byte(out.Choices[0].Delta.Content))
					}
					str := wordwrap.WrapString(out.Choices[0].Delta.Content, 25)
					fmt.Printf("\x1b[1:32m%s", str)
				})

			var event EventManager
			event.Errata(err)

			if !isTestingEnvironment() {
				bWriter.Write([]byte("\n\n###\n\n"))
				buffer = append(buffer, "\n\n###\n\n")
			}
			out := strings.Join(buffer, "")
			for i := range sresp.Choices {
				sresp.Choices[i].Delta.Content = fmt.Sprint(util.RemoveWrapper(out))
			}

			if !isTestingEnvironment() {
				node.layout.app.Sync()
			}

			c.chatStreamResponse = sresp
			return c.chatStreamResponse, nil
		}

		client := *service.client
		resp, err := client.ChatCompletion(
			service.ctx,
			req)

		var event EventManager
		event.Errata(err)

		if !isTestingEnvironment() {
			node.layout.app.Sync()
		}
		c.chatResponse = resp
		return nil, c.chatResponse

	}
	return nil, nil
}

// SendCompletionPrompt - Send task prompt on stream mode
func (c *Prompt) SendCompletionPrompt(service Agent) *gpt3.CompletionResponse {
	if isContextValid(service) {
		var buffer []string

		resp := &gpt3.CompletionResponse{}

		msg := service.SetTemplate(service.cachedPrompt, service.PromptProperties.Input[0])

		if !isTestingEnvironment() {
			service.PromptProperties.MaxTokens = 1024 + len(util.EncodePromptBytePair(msg, service.EngineProperties.Model))
			node.controller.currentAgent.preferences.MaxTokens = service.PromptProperties.MaxTokens
		}
		req := gpt3.CompletionRequest{
			Prompt:           msg,
			MaxTokens:        gpt3.IntPtr(service.PromptProperties.MaxTokens),
			Temperature:      gpt3.Float32Ptr(service.EngineProperties.Temperature),
			TopP:             gpt3.Float32Ptr(service.EngineProperties.TopP),
			PresencePenalty:  *gpt3.Float32Ptr(service.EngineProperties.PresencePenalty),
			FrequencyPenalty: *gpt3.Float32Ptr(service.EngineProperties.FrequencyPenalty),
			Stream:           service.preferences.IsPromptStreaming,
			N:                gpt3.IntPtr(service.PromptProperties.Results),
			LogProbs:         gpt3.IntPtr(service.PromptProperties.Probabilities),
			Echo:             false}

		if service.preferences.IsPromptStreaming {

			var bWriter tview.TextViewWriter
			if !isTestingEnvironment() {
				bWriter = node.layout.promptOutput.BatchWriter()
				defer bWriter.Close()
				bWriter.Clear()
			}
			fmt.Print("\033[H\033[2J")
			isOnce := false
			client := *service.client
			err := client.CompletionStreamWithEngine(
				service.ctx,
				service.EngineProperties.Model,
				req, func(out *gpt3.CompletionResponse) {
					go func(in chan string) {
						if !isOnce {
							resp.ID = out.ID
							resp.Choices = out.Choices
							resp.Created = out.Created
							resp.Model = out.Model
							resp.Object = out.Object
							resp.Usage = out.Usage
							isOnce = true
						}

						resp.Choices[0].FinishReason = out.Choices[0].FinishReason
						resp.Choices[0].Index = out.Choices[0].Index

						resp.Choices = append(resp.Choices, out.Choices[0])
						resp.Choices[0].LogProbs.TextOffset = append(resp.Choices[0].LogProbs.TextOffset, out.Choices[0].LogProbs.TextOffset...)
						resp.Choices[0].LogProbs.TokenLogprobs = append(resp.Choices[0].LogProbs.TokenLogprobs, out.Choices[0].LogProbs.TokenLogprobs...)
						resp.Choices[0].LogProbs.Tokens = append(resp.Choices[0].LogProbs.Tokens, out.Choices[0].LogProbs.Tokens...)
						resp.Choices[0].LogProbs.TopLogprobs = append(resp.Choices[0].LogProbs.TopLogprobs, out.Choices[0].LogProbs.TopLogprobs...)

						for i := range out.Choices {
							if !isTestingEnvironment() {
								buffer = append(buffer, out.Choices[i].Text)
							}
							in <- out.Choices[i].Text
							str := wordwrap.WrapString(out.Choices[i].Text, 50)
							fmt.Printf("\x1b[1:32m%s", str)
						}
					}(service.preferences.InlineText)
					if !isTestingEnvironment() {
						bWriter.Write([]byte(<-service.preferences.InlineText))
					}
				})

			var event EventManager
			event.Errata(err)

			out := strings.Join(buffer, "")
			for i := range resp.Choices {
				resp.Choices[i].Text = fmt.Sprint(util.RemoveWrapper(out))
			}

			if !isTestingEnvironment() {
				bWriter.Write([]byte("\n\n###\n\n"))
				node.layout.app.Sync()
			}
			c.contextualResponse = resp
			return c.contextualResponse
		}

		client := *service.client
		resp, err := client.CompletionWithEngine(
			service.ctx,
			service.EngineProperties.Model,
			req)

		var event EventManager
		event.Errata(err)

		if isTestingEnvironment() {
			node.layout.app.Sync()
		}
		c.contextualResponse = resp
		return c.contextualResponse

	}
	return nil
}

// SendEditPrompt - Send edit instruction task prompt
func (c *Prompt) SendEditPrompt(service Agent) *gpt3.EditsResponse {
	if isContextValid(service) &&
		service.PromptProperties.Input != nil {

		req := gpt3.EditsRequest{
			Model:       service.EngineProperties.Model,
			Input:       service.PromptProperties.Input[0],
			Instruction: service.PromptProperties.Instruction[0],
			Temperature: gpt3.Float32Ptr(service.EngineProperties.Temperature),
			TopP:        gpt3.Float32Ptr(service.EngineProperties.TopP),
			N:           gpt3.IntPtr(service.PromptProperties.Results)}

		client := *service.client
		resp, err := client.Edits(
			service.ctx,
			req)

		var event EventManager
		event.Errata(err)

		if !isTestingEnvironment() {
			node.layout.app.Sync()
		}

		c.extendedResponse = resp
		return c.extendedResponse
	}
	return nil
}

// SendEmbeddingPrompt - Creates an embedding vector representing the input text
func (c *Prompt) SendEmbeddingPrompt(service Agent) *gpt3.EmbeddingsResponse {
	if isContextValid(service) {
		req := gpt3.EmbeddingsRequest{
			Model: service.EngineProperties.Model,
			Input: service.PromptProperties.Input,
		}

		client := *service.client
		resp, err := client.Embeddings(
			service.ctx,
			req)

		var event EventManager
		event.Errata(err)

		if !isTestingEnvironment() {
			node.layout.app.Sync()
		}

		c.embeddingResponse = resp
		return c.embeddingResponse
	}
	return nil
}

// SendPredictablePrompt - Send a predictable request
func (c *Prompt) SendPredictablePrompt(service Agent) *model.PredictResponse {
	isValid := isContextValid(service)
	if isValid {
		req := model.PredictRequest{
			Document: string(service.PredictProperties.Input[0]),
		}

		var event EventManager
		out, err := json.Marshal(req)
		if err != nil {
			event.Errata(err)
			return nil
		}

		var body io.Reader = bytes.NewBuffer(out)

		path := string("https://api.gptzero.me/v2/predict/text")

		if out != nil {
			req, err := http.NewRequestWithContext(service.ctx, "POST", path, body)
			if err != nil {
				println(err)
				event.Errata(err)
				return nil
			}

			if req != nil {
				req.Header.Set("Accept", "application/json")
				req.Header.Set("Content-type", "application/json")
				req.Header.Set("X-Api-Key", service.key[1])

				resp, err := service.exClient.Do(req)
				if err != nil {
					event.Errata(err)
					return nil
				}

				if resp != nil {
					defer resp.Body.Close()
					data, _ := io.ReadAll(resp.Body)
					var dataReader io.Reader = bytes.NewBuffer(data)

					in := new(model.PredictResponse)
					json.NewDecoder(dataReader).Decode(in)

					event.Errata(err)

					c.predictableResponse = in
					return c.predictableResponse
				}
			}
		}
	}
	return nil
}

// GetListModels - Get actual list of available models
func (c *Prompt) GetListModels(service Agent) *gpt3.EnginesResponse {
	client := *service.client
	resp, err := client.Engines(service.ctx)

	var event EventManager
	event.Errata(err)

	return resp
}
