// Package service section
package service

import (
	"context"
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
