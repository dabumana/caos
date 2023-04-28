// Package service section
package service

import (
	"context"
	"fmt"
	"os"

	"caos/model"
	"caos/service/parameters"
	"caos/util"
	"encoding/csv"
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Agent - Contextual client API
type Agent struct {
	// Versioned and ID
	version string
	id      string
	key     []string
	// Assistant context
	templateID  []string
	templateCtx []string
	// Context
	ctx context.Context
	// Client
	client   *gpt3.Client
	exClient *http.Client
	// Properties
	engineProperties  model.EngineProperties
	promptProperties  model.PromptProperties
	predictProperties model.PredictProperties
	// Preferences
	preferences parameters.GlobalPreferences
	// Temporal cache
	cachedPrompt string
}

// Initialize - Creates context background to be used along with the client
func (c *Agent) Initialize() Agent {
	// ID
	c.version = "v.0.2.2"
	c.id = "anon"
	// Key
	c.key = getKeys()
	// template
	c.templateID, c.templateCtx = getTemplateFromLocal()
	// Background context
	c.ctx = context.Background()
	c.client, c.exClient = c.Connect()
	// Role
	c.preferences.Role = model.Assistant
	// Global preferences
	c.preferences.TemplateIDs = len(c.templateID)
	c.preferences.Template = 0
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
	c.preferences.IsPromptStreaming = true
	c.preferences.IsTurbo = false
	c.preferences.InlineText = make(chan string)
	// Return created client
	return *c
}

// Connect - Contextualize the API to create a new client
func (c *Agent) Connect() (*gpt3.Client, *http.Client) {
	godotenv.Load()

	externalClient := http.Client{
		Transport: http.DefaultTransport,
	}

	option := gpt3.WithHTTPClient(&externalClient)
	client := gpt3.NewClient(c.key[0], option)

	c.client = &client
	c.exClient = &externalClient

	return c.client, c.exClient
}

// SaveKeys - Set API keys
func (c *Agent) SaveKeys() {
	event := &EventManager{}

	dir, _ := os.Getwd()
	path := fmt.Sprint(dir, ".env")

	_, err := os.Open(path)
	if err != nil {
		file, err := os.Create(".env")
		if err != nil {
			event.Errata(err)
		}
		outFile := fmt.Sprintf("API_KEY=%v\nZERO_API_KEY=%v\n", c.key[0], c.key[1])
		file.Write([]byte(outFile))
		file.Sync()
	}
}

// GetStatus - Current agent information
func (c *Agent) GetStatus() parameters.GlobalPreferences {
	return c.preferences
}

// getKeys - Grab API keys
func getKeys() []string {
	dir, _ := os.Getwd()
	path := fmt.Sprintf("%v/.env", dir)

	file, _ := os.Stat(path)
	if file != nil {
		return getKeyFromLocal()
	}
	return getKeyFromEnv()
}

// getKeyFromEnv - Get environment keys
func getKeyFromEnv() []string {
	var keys []string
	// Variables from environment
	api := os.Getenv("API_KEY")
	if api != "" {
		keys = append(keys, api)
	} else {
		keys = append(keys, "")
	}

	zeroAPI := os.Getenv("ZERO_API_KEY")
	if zeroAPI != "" {
		keys = append(keys, zeroAPI)
	} else {
		keys = append(keys, "")
	}

	return keys
}

// getKeyFromLocal - Get the currect key stablished on the environment
func getKeyFromLocal() []string {
	var keys []string

	dir, _ := os.Getwd()
	path := fmt.Sprintf("%v/.env", dir)

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		keys = append(keys, "", "")
		return keys
	}

	val1, _ := viper.Get("API_KEY").(string)
	val2, _ := viper.Get("ZERO_API_KEY").(string)

	keys = append(keys, val1, val2)

	return keys
}

// getTemplateFromLocal - Get templates on local dir
func getTemplateFromLocal() ([]string, []string) {
	var index []string
	var context []string

	dir, _ := os.Getwd()
	path := fmt.Sprintf("%v/template/role.csv", dir)

	file, _ := os.Open(path)
	reader := csv.NewReader(file)
	data, _ := reader.ReadAll()

	for _, j := range data {
		for k, l := range j {
			if k == 0 {
				index = append(index, l)
			} else if k == 1 {
				context = append(context, l)
			}
		}
	}

	return index, context
}

// SetEngineParameters - Set engine parameters for the current prompt
func (c *Agent) SetEngineParameters(id string, pmodel string, role model.Roles, temperature float32, topp float32, penalty float32, frequency float32) model.EngineProperties {
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
func (c *Agent) SetPromptParameters(promptContext []string, instruction []string, tokens int, results int, probabilities int) model.PromptProperties {
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
func (c *Agent) SetPredictionParameters(prompContext []string) model.PredictProperties {
	properties := model.PredictProperties{
		Input: prompContext,
	}
	return properties
}

// SetPrompt - Conversion human-ai roles
func (c *Agent) SetPrompt(context string, input string) []string {
	if context == "" {
		context = c.templateCtx[c.preferences.Template]
	}

	prompt := []string{context + input}
	return prompt
}
