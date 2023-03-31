// Package service section
package service

import (
	"bytes"
	"caos/model"
	"caos/service/parameters"
	"caos/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
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

// IsContextValid - Client context validation
func IsContextValid(current Agent) bool {
	if current.ctx == nil {
		log.Fatalln("Context NOT found")
		return false
	} else if current.client == nil {
		log.Fatalln("Client NOT found")
		return false
	}

	return true
}

// SendStreamingChatCompletion - Send streaming chat completion prompt
func (c Prompt) SendStreamingChatCompletion(service Agent) *gpt3.ChatCompletionStreamResponse {
	if IsContextValid(service) {
		var buffer []string
		var event EventManager

		resp := &gpt3.ChatCompletionStreamResponse{}

		msg := gpt3.ChatCompletionRequestMessage{
			Role:    string(service.preferences.Role),
			Content: node.controller.currentAgent.SetPrompt(service.promptProperties.PromptContext[0])[0],
		}

		req := gpt3.ChatCompletionRequest{
			Model:            service.engineProperties.Model,
			User:             service.id,
			Messages:         []gpt3.ChatCompletionRequestMessage{msg},
			MaxTokens:        *gpt3.IntPtr(service.promptProperties.MaxTokens),
			Temperature:      *gpt3.Float32Ptr(service.engineProperties.Temperature),
			TopP:             *gpt3.Float32Ptr(service.engineProperties.TopP),
			PresencePenalty:  *gpt3.Float32Ptr(service.engineProperties.PresencePenalty),
			FrequencyPenalty: *gpt3.Float32Ptr(service.engineProperties.FrequencyPenalty),
			Stream:           false,
			N:                *gpt3.IntPtr(service.promptProperties.Results),
		}

		bWriter := node.layout.promptOutput.BatchWriter()
		defer bWriter.Close()
		bWriter.Clear()

		bWriter.Write([]byte("\n"))
		buffer = append(buffer, "\n")

		fmt.Print("\033[H\033[2J")
		err := service.client.ChatCompletionStream(
			node.controller.currentAgent.ctx,
			req, func(out *gpt3.ChatCompletionStreamResponse) {
				resp.ID = out.ID
				resp.Choices = out.Choices
				resp.Created = out.Created
				resp.Model = out.Model
				resp.Object = out.Object
				resp.Usage = out.Usage
				// Choices
				resp.Choices[0].Index = out.Choices[0].Index
				resp.Choices[0].FinishReason = out.Choices[0].FinishReason
				// Delta
				resp.Choices[0].Delta = out.Choices[0].Delta
				resp.Choices[0].Delta.Content = out.Choices[0].Delta.Content
				resp.Choices[0].Delta.Role = out.Choices[0].Delta.Role
				// Write buffer
				buffer = append(buffer, out.Choices[0].Delta.Content)
				bWriter.Write([]byte(out.Choices[0].Delta.Content))
				if out.Choices[0].Delta.Content == "\n" {
					fmt.Print("\r" + out.Choices[0].Delta.Content)
				} else {
					fmt.Print(out.Choices[0].Delta.Content)
				}
			})

		event.Errata(err)

		bWriter.Write([]byte("\n\n###\n\n"))
		buffer = append(buffer, "\n\n###\n\n")

		out := strings.Join(buffer, "")
		resp.Choices[0].Delta.Content = fmt.Sprint(util.RemoveWrapper(out))

		node.layout.app.Sync()
		c.chatStreamResponse = resp
		return c.chatStreamResponse
	}
	return nil
}

// SendChatCompletion - Send chat completion prompt
func (c Prompt) SendChatCompletion(service Agent) *gpt3.ChatCompletionResponse {
	if IsContextValid(service) {
		msg := gpt3.ChatCompletionRequestMessage{
			Role:    string(service.preferences.Role),
			Content: node.controller.currentAgent.SetPrompt(service.promptProperties.PromptContext[0])[0],
		}

		req := gpt3.ChatCompletionRequest{
			Model:            service.engineProperties.Model,
			User:             service.id,
			Messages:         []gpt3.ChatCompletionRequestMessage{msg},
			MaxTokens:        *gpt3.IntPtr(service.promptProperties.MaxTokens),
			Temperature:      *gpt3.Float32Ptr(service.engineProperties.Temperature),
			TopP:             *gpt3.Float32Ptr(service.engineProperties.TopP),
			PresencePenalty:  *gpt3.Float32Ptr(service.engineProperties.PresencePenalty),
			FrequencyPenalty: *gpt3.Float32Ptr(service.engineProperties.FrequencyPenalty),
			Stream:           false,
			N:                *gpt3.IntPtr(service.promptProperties.Results),
		}

		resp, err := service.client.ChatCompletion(
			node.controller.currentAgent.ctx,
			req)

		var event EventManager
		event.Errata(err)

		node.layout.app.Sync()
		c.chatResponse = resp
		return c.chatResponse
	}
	return nil
}

// SendCompletion - Send task prompt
func (c Prompt) SendCompletion(service Agent) *gpt3.CompletionResponse {
	if IsContextValid(service) {
		req := gpt3.CompletionRequest{
			Prompt:           node.controller.currentAgent.SetPrompt(service.promptProperties.PromptContext[0]),
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
			node.controller.currentAgent.ctx,
			service.engineProperties.Model,
			req)

		var event EventManager
		event.Errata(err)

		node.layout.app.Sync()
		c.contextualResponse = resp
		return c.contextualResponse
	}
	return nil
}

// SendStreamingCompletion - Send task prompt on stream mode
func (c Prompt) SendStreamingCompletion(service Agent) *gpt3.CompletionResponse {
	if IsContextValid(service) {
		var event EventManager
		var buffer []string

		resp := &gpt3.CompletionResponse{}

		req := gpt3.CompletionRequest{
			Prompt:           node.controller.currentAgent.SetPrompt(service.promptProperties.PromptContext[0]),
			MaxTokens:        gpt3.IntPtr(service.promptProperties.MaxTokens),
			Temperature:      gpt3.Float32Ptr(service.engineProperties.Temperature),
			TopP:             gpt3.Float32Ptr(service.engineProperties.TopP),
			PresencePenalty:  *gpt3.Float32Ptr(service.engineProperties.PresencePenalty),
			FrequencyPenalty: *gpt3.Float32Ptr(service.engineProperties.FrequencyPenalty),
			Stream:           false,
			N:                gpt3.IntPtr(service.promptProperties.Results),
			LogProbs:         gpt3.IntPtr(service.promptProperties.Probabilities),
			Echo:             true}

		bWriter := node.layout.promptOutput.BatchWriter()
		defer bWriter.Close()
		bWriter.Clear()

		bWriter.Write([]byte("\n"))
		buffer = append(buffer, "\n")

		fmt.Print("\033[H\033[2J")
		isOnce := false
		err := service.client.CompletionStreamWithEngine(
			node.controller.currentAgent.ctx,
			service.engineProperties.Model,
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
						buffer = append(buffer, out.Choices[i].Text)
						in <- out.Choices[i].Text
						if out.Choices[i].Text == "\n" {
							fmt.Print("\r" + out.Choices[i].Text)
						} else {
							fmt.Print(out.Choices[i].Text)
						}
					}
				}(node.controller.currentAgent.preferences.InlineText)
				// Write buffer
				bWriter.Write([]byte(<-node.controller.currentAgent.preferences.InlineText))
			})

		event.Errata(err)

		bWriter.Write([]byte("\n\n###\n\n"))
		buffer = append(buffer, "\n\n###\n\n")

		resp.Choices[0].Text = fmt.Sprint(buffer)

		node.layout.app.Sync()
		c.contextualResponse = resp
		return c.contextualResponse
	}
	return nil
}

// SendEditPrompt - Send edit instruction task prompt
func (c Prompt) SendEditPrompt(service Agent) *gpt3.EditsResponse {
	if IsContextValid(service) && service.promptProperties.PromptContext != nil {
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

		node.layout.app.Sync()
		c.extendedResponse = resp
		return c.extendedResponse
	}
	return nil
}

// SendEmbeddingPrompt - Creates an embedding vector representing the input text
func (c Prompt) SendEmbeddingPrompt(service Agent) *gpt3.EmbeddingsResponse {
	if IsContextValid(service) {
		req := gpt3.EmbeddingsRequest{
			Model: service.engineProperties.Model,
			Input: service.promptProperties.PromptContext,
		}

		resp, err := service.client.Embeddings(
			service.ctx,
			req)

		var event EventManager
		event.Errata(err)

		node.layout.app.Sync()
		c.embeddingResponse = resp
		return c.embeddingResponse
	}
	return nil
}

// SendPredictablePrompt - Send a predictable request
func (c Prompt) SendPredictablePrompt(service Agent) *model.PredictResponse {
	isValid := IsContextValid(service)
	if isValid {
		req := model.PredictRequest{
			Document: string(service.predictProperties.Input[0]),
		}

		zeroApiKey := os.Getenv("ZERO_API_KEY")
		if zeroApiKey == "" {
			log.Fatalln("Missing GPT ZERO API KEY, you have some free requests.")
		}

		var event EventManager
		out, err := json.Marshal(req)
		if err != nil {
			event.Errata(err)
			return nil
		}

		var body io.Reader = bytes.NewBuffer(out)
		path := parameters.ExternalBaseURL + string("/v2/predict/text")

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
				req.Header.Set("X-Api-Key", zeroApiKey)

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
func (c Prompt) GetListModels(service Agent) *gpt3.EnginesResponse {
	resp, err := service.client.Engines(service.ctx)

	var event EventManager
	event.Errata(err)

	return resp
}
