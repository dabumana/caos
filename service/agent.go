// Package service section
package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"caos/model"
	"caos/service/parameters"
	"caos/util"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Agent - Contextual client API
type Agent struct {
	id                string
	keys              []string
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
	// ID
	c.id = "anon"
	// Key
	c.keys = c.getKeyFromLocalEnv()
	// Role
	c.preferences.Role = model.Assistant
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
	c.preferences.IsDeveloper = false
	c.preferences.IsEditable = false
	c.preferences.IsLoading = false
	c.preferences.IsNewSession = true
	c.preferences.IsPromptReady = false
	c.preferences.IsTraining = false
	c.preferences.IsPromptStreaming = true
	c.preferences.IsTurbo = false
	c.preferences.InlineText = make(chan string)
	// Template
	reader, _ := ioutil.ReadDir("../template/")
	for _, file := range reader {
		c.preferences.Template = append(c.preferences.Template, file.Name())
	}
	// Return created client
	return c
}

// Connect - Contextualize the API to create a new client
func (c Agent) Connect() (gpt3.Client, *http.Client) {
	godotenv.Load()

	externalClient := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Duration(12 * time.Second),
	}

	option := gpt3.WithHTTPClient(&externalClient)
	client := gpt3.NewClient(c.keys[0], option)

	c.client = client
	c.exClient = &externalClient

	return c.client, c.exClient
}

// getKeyFromLocalEnv - Get the currect key stablished on the environment
func (c Agent) getKeyFromLocalEnv() []string {
	var keys []string

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}

	val1, _ := viper.Get("API_KEY").(string)
	val2, _ := viper.Get("ZERO_API_KEY").(string)

	if val1 != "" {
		keys = append(keys, val1)
	}
	if val2 != "" {
		keys = append(keys, val2)
	}

	return keys
}

// SetEngineParameters - Set engine parameters for the current prompt
func (c Agent) SetEngineParameters(id string, pmodel string, role model.Roles, temperature float32, topp float32, penalty float32, frequency float32) model.EngineProperties {
	properties := model.EngineProperties{
		UserID:           id,
		Model:            pmodel,
		Role:             role,
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

// SetPrompt - Conversion human-ai roles
func (c Agent) SetPrompt(context string) []string {
	var prompt []string
	if node.controller.currentAgent.preferences.IsConversational &&
		!node.controller.currentAgent.preferences.IsDeveloper {
		prompt = []string{fmt.Sprintf("Human: %v \nAI:", context)}
	} else if node.controller.currentAgent.preferences.IsDeveloper &&
		!node.controller.currentAgent.preferences.IsConversational {
		prompt = []string{fmt.Sprintf("Developer Mode: %v \nAI:", context)}
	} else {
		prompt = []string{context}
	}
	return prompt
}
