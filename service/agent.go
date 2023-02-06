// Package service section
package service

import (
	"context"
	"log"
	"net/http"
	"os"

	"caos/model"
	"caos/service/parameters"
	"caos/util"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

// Agent - Contextual client API
type Agent struct {
	ctx               context.Context
	client            gpt3.Client
	exClient          *http.Client
	engineProperties  model.EngineProperties
	promptProperties  model.PromptProperties
	predictProperties model.PredictProperties
	preferences       parameters.GlobalPreferences
}

// Initialize - Creates context background to be used along with the client
func (c Agent) Initialize() Agent {
	// Background context
	c.ctx = context.Background()
	c.client, c.exClient = c.Connect()
	// Global preferences
	c.preferences.Engine = "text-davinci-003"
	c.preferences.Frequency = util.ParseFloat32("\u0030\u002e\u0035")
	c.preferences.Penalty = util.ParseFloat32("\u0030\u002e\u0035")
	c.preferences.MaxTokens = util.ParseInt64("\u0032\u0035\u0030")
	c.preferences.Mode = "Text"
	c.preferences.Models = append(c.preferences.Models, "zero-gpt")
	c.preferences.Probabilities = util.ParseInt32("\u0031")
	c.preferences.Results = util.ParseInt32("\u0031")
	c.preferences.Temperature = util.ParseFloat32("\u0030\u002e\u0034")
	c.preferences.Topp = util.ParseFloat32("\u0030\u002e\u0036")
	// Mode selection
	c.preferences.IsConversational = false
	c.preferences.IsEditable = false
	c.preferences.IsLoading = false
	c.preferences.IsNewSession = true
	c.preferences.IsPredictable = false
	c.preferences.IsPromptReady = false
	c.preferences.IsTraining = false
	// Return created client
	return c
}

// Connect - Contextualize the API to create a new client
func (c Agent) Connect() (gpt3.Client, *http.Client) {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	client := gpt3.NewClient(apiKey)
	externalClient := http.Client{
		Transport: http.DefaultTransport,
	}

	c.client = client
	c.exClient = &externalClient

	return c.client, c.exClient
}

// SetEngineParameters - Set engine parameters for the current prompt
func (c Agent) SetEngineParameters(pmodel string, temperature float32, topp float32, penalty float32, frequency float32) model.EngineProperties {
	properties := model.EngineProperties{
		Model:            pmodel,
		Temperature:      temperature,
		TopP:             topp,
		PresencePenalty:  penalty,
		FrequencyPenalty: frequency,
	}
	return properties
}

// SetPromptParameters - Set request parameters for the current prompt
func (c Agent) SetPromptParameters(promptContext []string, instruction []string, tokens int, results int, probabilities int) model.PromptProperties {
	properties := model.PromptProperties{
		PromptContext: promptContext,
		Instruction:   instruction,
		MaxTokens:     tokens,
		Results:       results,
		Probabilities: probabilities,
	}
	return properties
}

// SetPredictionParameters - Set prediction parameters for the current prompt
func (c Agent) SetPredictionParameters(prompContext []string) model.PredictProperties {
	properties := model.PredictProperties{
		Input: prompContext,
	}
	return properties
}
