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

// Client service - Use the clients properties based on the parental role
type ServiceClient struct {
	ctx              context.Context
	client           gpt3.Client
	engineProperties model.EngineProperties
	promptProperties model.PromptProperties
}

/* Service client functionality*/
// Intialize
func (c ServiceClient) Initialize() ServiceClient {
	c.ctx = context.Background()
	c.client = c.Connect()
	return c
}

// Connect
func (c ServiceClient) Connect() gpt3.Client {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	client := gpt3.NewClient(apiKey)
	c.client = client
	return c.client
}

// Log client context
func (c ServiceClient) LogClient() {
	fmt.Printf("-------------------------------------------\n")
	fmt.Printf("Context: %v\nClient: %v\n", c.ctx, c.client)
	fmt.Printf("-------------------------------------------\n")
}

// Log current engine
func (c ServiceClient) LogEngine() {
	Node.Layout.metadataOutput.SetText(
		fmt.Sprintf("\nModel: %v\nTemperature: %v\nTopp: %v\nFrequency penalty: %v\nPresence penalty: %v\nPrompt context: %v\nPrompt: %v\nProbabilities: %v\nResults: %v\nMax tokens: %v\n",
			c.engineProperties.Model,
			c.engineProperties.Temperature,
			c.engineProperties.TopP,
			c.engineProperties.FrequencyPenalty,
			c.engineProperties.PresencePenalty,
			c.promptProperties.PromptContext,
			c.promptProperties.Prompt,
			c.promptProperties.Probabilities,
			c.promptProperties.Results,
			c.promptProperties.MaxTokens))
}

// Set engine parameters
func (c ServiceClient) SetEngineParameters(pmodel string, temperature float32, topp float32, penalty float32, frequency float32) model.EngineProperties {
	engine := model.EngineProperties{
		Model:            pmodel,
		Temperature:      temperature,
		TopP:             topp,
		PresencePenalty:  penalty,
		FrequencyPenalty: frequency,
	}
	return engine
}

// Set request parameters
func (c ServiceClient) SetRequestParameters(promptContext []string, prompt []string, tokens int, results int, probabilities int) model.PromptProperties {
	request := model.PromptProperties{
		PromptContext: promptContext,
		Prompt:        prompt,
		MaxTokens:     tokens,
		Results:       results,
		Probabilities: probabilities,
	}
	return request
}
