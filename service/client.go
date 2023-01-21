// Package service section
package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"caos/model"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

// Client - Contextual client API
type Client struct {
	ctx              context.Context
	client           gpt3.Client
	engineProperties model.EngineProperties
	promptProperties model.PromptProperties
}

// Initialize - Creates context background to be used along with the client
func (c Client) Initialize() Client {
	c.ctx = context.Background()
	c.client = c.Connect()
	return c
}

// Connect - Contextualize the API to create a new client
func (c Client) Connect() gpt3.Client {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	client := gpt3.NewClient(apiKey)
	c.client = client
	return c.client
}

// LogClient - Log client context
func (c Client) LogClient() {
	fmt.Printf("-------------------------------------------\n")
	fmt.Printf("Context: %v\nClient: %v\n", c.ctx, c.client)
	fmt.Printf("-------------------------------------------\n")
}

// LogEngine - Log current engine
func (c Client) LogEngine() {
	node.layout.metadataOutput.SetText(
		fmt.Sprintf("\nModel: %v\nTemperature: %v\nTopp: %v\nFrequency penalty: %v\nPresence penalty: %v\nPrompt: %v\nInstruction: %v\nProbabilities: %v\nResults: %v\nMax tokens: %v\n",
			c.engineProperties.Model,
			c.engineProperties.Temperature,
			c.engineProperties.TopP,
			c.engineProperties.FrequencyPenalty,
			c.engineProperties.PresencePenalty,
			c.promptProperties.PromptContext,
			c.promptProperties.Instruction,
			c.promptProperties.Probabilities,
			c.promptProperties.Results,
			c.promptProperties.MaxTokens))
}

// SetEngineParameters - Set engine parameters for the current prompt
func (c Client) SetEngineParameters(pmodel string, temperature float32, topp float32, penalty float32, frequency float32) model.EngineProperties {
	engine := model.EngineProperties{
		Model:            pmodel,
		Temperature:      temperature,
		TopP:             topp,
		PresencePenalty:  penalty,
		FrequencyPenalty: frequency,
	}
	return engine
}

// SetRequestParameters - Set request parameters for the current prompt
func (c Client) SetRequestParameters(promptContext []string, prompt []string, tokens int, results int, probabilities int) model.PromptProperties {
	request := model.PromptProperties{
		PromptContext: promptContext,
		Instruction:   prompt,
		MaxTokens:     tokens,
		Results:       results,
		Probabilities: probabilities,
	}
	return request
}
