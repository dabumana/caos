// Package service section
package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"caos/model"
	"caos/service/parameters"
	"caos/util"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/gdamore/tcell/v2"
)

// EventPool - Historical events
var EventPool []model.HistoricalEvent

// SessionPool - Historical sessions
var SessionPool []model.HistoricalSession

// TrainingEventPool - Historical training events
var TrainingEventPool []model.HistoricalTrainingEvent

// TrainingSessionPool - Historical training sessions
var TrainingSessionPool []model.HistoricalTrainingSession

// CurrentID - contextual parent id
var CurrentID string

// EventManager - Log event service
type EventManager struct {
	event   model.HistoricalEvent
	session model.HistoricalSession
}

// SaveLog - Save log with actual historic detail
func (c EventManager) SaveLog() {
	raw, _ := json.MarshalIndent(SessionPool[len(SessionPool)-1], "", "\u0009")
	out := util.ConstructPathFileTo("log", "json")
	out.WriteString(string(raw))
}

// ClearSession - Clear all the pools
func (c EventManager) ClearSession() {
	EventPool = nil
	SessionPool = nil
	TrainingEventPool = nil
	TrainingSessionPool = nil
}

// AppendToSession - Add a set of events as a session
func (c EventManager) AppendToSession(header *model.EngineProperties, body *model.PromptProperties, id string, train model.TrainingPrompt) {
	c.event.Event.Header = *header
	c.event.Event.Body = *body
	c.event.Timestamp = fmt.Sprint(time.Now().UnixMilli())

	EventPool = append(EventPool, c.event)

	c.session.ID = id
	c.session.Session = EventPool

	if parameters.IsTraining {

		event := model.HistoricalTrainingEvent{
			Timestamp: c.event.Timestamp,
			Event:     train,
		}

		TrainingEventPool = append(TrainingEventPool, event)

		session := model.HistoricalTrainingSession{
			ID:      c.session.ID,
			Session: []model.HistoricalTrainingEvent{event},
		}

		TrainingSessionPool = append(TrainingSessionPool, session)
	}

	SessionPool = append(SessionPool, c.session)
}

// AppendToLayout - Append and visualize content in console page view
func (c EventManager) AppendToLayout(responses []string) {
	parameters.PromptCtx = responses
	log := strings.Join(responses, "")
	reg := strings.ReplaceAll(log, "[]", "\n")
	node.layout.promptOutput.SetText(reg)
}

// LogCompletion - Response details in a .json file
func (c EventManager) LogCompletion(header *model.EngineProperties, body *model.PromptProperties, resp *gpt3.CompletionResponse) {
	if parameters.IsNewSession {
		c.ClearSession()
		parameters.IsNewSession = false
	}

	modelTrainer := model.TrainingPrompt{
		Prompt:     body.PromptContext,
		Completion: []string{resp.Choices[0].Text},
	}

	c.AppendToSession(header, body, resp.ID, modelTrainer)

	CurrentID = resp.ID
}

// LogInstruction - Response details in a .json file
func (c EventManager) LogInstruction(header *model.EngineProperties, body *model.PromptProperties, resp *gpt3.EditsResponse) {
	modelTrainer := model.TrainingPrompt{
		Prompt:     body.PromptContext,
		Completion: []string{resp.Choices[0].Text},
	}

	c.AppendToSession(header, body, CurrentID, modelTrainer)
}

// VisualLogCompletion - Response details
func (c EventManager) VisualLogCompletion(resp *gpt3.CompletionResponse) {
	var responses []string
	for i := range resp.Choices {
		responses = append(responses, resp.Choices[i].Text, "\n\n###\n\n")
	}

	c.AppendToLayout(responses)

	node.layout.infoOutput.SetText(
		fmt.Sprintf("\nID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nToken probs: %v \nToken top: %v\n",
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

// VisualLogInstruction - Log edited response details
func (c EventManager) VisualLogInstruction(resp *gpt3.EditsResponse) {
	var responses []string
	for i := range resp.Choices {
		responses = append(responses, resp.Choices[i].Text, "\n\n###\n\n")
	}

	c.AppendToLayout(responses)

	node.layout.infoOutput.SetText(fmt.Sprintf("\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
		resp.Created,
		resp.Object,
		resp.Usage.CompletionTokens,
		resp.Usage.PromptTokens,
		resp.Usage.TotalTokens,
		resp.Choices[0].Index))
}

// LogClient - Log client context
func (c EventManager) LogClient(client Client) {
	fmt.Printf("-------------------------------------------\n")
	fmt.Printf("Context: %v\nClient: %v\n", client.ctx, client.client)
	fmt.Printf("-------------------------------------------\n")
}

// LogEngine - Log current engine
func (c EventManager) LogEngine(client Client) {
	node.layout.metadataOutput.SetText(
		fmt.Sprintf("\nModel: %v\nTemperature: %v\nTopp: %v\nFrequency penalty: %v\nPresence penalty: %v\nPrompt: %v\nInstruction: %v\nProbabilities: %v\nResults: %v\nMax tokens: %v\n",
			client.engineProperties.Model,
			client.engineProperties.Temperature,
			client.engineProperties.TopP,
			client.engineProperties.FrequencyPenalty,
			client.engineProperties.PresencePenalty,
			client.promptProperties.PromptContext,
			client.promptProperties.Instruction,
			client.promptProperties.Probabilities,
			client.promptProperties.Results,
			client.promptProperties.MaxTokens))
}

// Errata - Generic error method
func (c EventManager) Errata(err error) {
	if err != nil {
		parameters.IsNewSession = true
		node.layout.infoOutput.SetText(err.Error())
		node.layout.promptInput.SetPlaceholder("Press ENTER again to repeat the request.")
		node.layout.promptInput.SetPlaceholderTextColor(tcell.ColorDarkOrange)
	} else {
		node.layout.promptInput.SetPlaceholder("Type here...")
		node.layout.promptInput.SetPlaceholderTextColor(tcell.ColorBlack)
	}

	parameters.IsLoading = false
	node.layout.promptInput.SetText("")
}
