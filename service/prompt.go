// Package service section
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"caos/model"
	"caos/service/parameters"

	"github.com/PullRequestInc/go-gpt3"
)

// Prompt - Handle prompt request
type Prompt struct {
	contextualResponse  *gpt3.CompletionResponse
	extendedResponse    *gpt3.EditsResponse
	embeddingResponse   *gpt3.EmbeddingsResponse
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

// SendCompletion - Send task prompt
func (c Prompt) SendCompletion(service Agent) *gpt3.CompletionResponse {
	isValid := IsContextValid(service)
	if isValid {
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
			node.agent.currentAgent.ctx,
			service.engineProperties.Model,
			req)

		var event EventManager
		event.Errata(err)

		c.contextualResponse = resp
		return c.contextualResponse
	}
	return nil
}

// SendEditPrompt - Send edit instruction task prompt
func (c Prompt) SendEditPrompt(service Agent) *gpt3.EditsResponse {
	isValid := IsContextValid(service)
	if isValid && service.promptProperties.PromptContext != nil {
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
	return nil
}

// SendEmbeddingPrompt - Creates an embedding vector representing the input text
func (c Prompt) SendEmbeddingPrompt(service Agent) *gpt3.EmbeddingsResponse {
	isValid := IsContextValid(service)
	if isValid {
		req := gpt3.EmbeddingsRequest{
			Model: service.engineProperties.Model,
			Input: service.promptProperties.PromptContext,
		}

		resp, err := service.client.Embeddings(
			service.ctx,
			req)

		var event EventManager
		event.Errata(err)

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

		var event EventManager
		out, err := json.Marshal(req)
		if err != nil {
			CleanConsoleView()
			event.Errata(err)
			return nil
		}

		var body io.Reader = bytes.NewBuffer(out)
		path := parameters.ExternalBaseURL + string("/v2/predict/text")

		if out != nil {
			req, err := http.NewRequestWithContext(service.ctx, "POST", path, body)
			if err != nil {
				println(err)
				CleanConsoleView()
				event.Errata(err)
				return nil
			}

			if req != nil {
				req.Header.Set("Content-type", "application/json")

				resp, err := service.exClient.Do(req)
				if err != nil {
					CleanConsoleView()
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
