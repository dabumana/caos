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
	out := util.ConstructPathFileToJSON("log")
	out.WriteString(string(raw))
}

// ClearSession - Clear all the pools
func (c EventManager) ClearSession() {
	EventPool = nil
	SessionPool = nil
	TrainingEventPool = nil
	TrainingSessionPool = nil
}

// Log - Response details in a .json file
func (c EventManager) Log(header *model.EngineProperties, body *model.PromptProperties, resp *gpt3.CompletionResponse) {
	if parameters.IsNewSession {
		c.ClearSession()
		parameters.IsNewSession = false
	}

	c.event.Event.Header = *header
	c.event.Event.Body = *body
	c.event.Timestamp = fmt.Sprint(time.Now().UnixMilli())

	EventPool = append(EventPool, c.event)

	c.session.ID = resp.ID
	c.session.Session = EventPool

	if parameters.IsTraining {
		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{resp.Choices[0].Text},
		}

		event := model.HistoricalTrainingEvent{
			Timestamp: c.event.Timestamp,
			Event:     modelTrainer,
		}

		TrainingEventPool = append(TrainingEventPool, event)

		session := model.HistoricalTrainingSession{
			ID:      c.session.ID,
			Session: []model.HistoricalTrainingEvent{event},
		}

		TrainingSessionPool = append(TrainingSessionPool, session)
	}

	SessionPool = append(SessionPool, c.session)

	CurrentID = c.session.ID
}

// LogEdit - Response details in a .json file
func (c EventManager) LogEdit(header *model.EngineProperties, body *model.PromptProperties, resp *gpt3.EditsResponse) {
	c.event.Event.Header = *header
	c.event.Event.Body = *body
	c.event.Timestamp = fmt.Sprint(time.Now().UnixMilli())

	EventPool = append(EventPool, c.event)

	c.session.ID = CurrentID
	c.session.Session = EventPool

	if parameters.IsTraining {
		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{resp.Choices[0].Text},
		}

		event := model.HistoricalTrainingEvent{
			Timestamp: c.event.Timestamp,
			Event:     modelTrainer,
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

// LogViz - Response details
func (c EventManager) LogViz(resp *gpt3.CompletionResponse) {
	var responses []string
	for i := range resp.Choices {
		responses = append(responses, resp.Choices[i].Text, "\n\n###\n\n")
	}
	parameters.PromptCtx = responses
	log := strings.Join(responses, "")
	reg := strings.ReplaceAll(log, "[]", "\n")
	node.layout.promptOutput.SetText(reg)
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

// LogVizEdit - Log edited response details
func (c EventManager) LogVizEdit(resp *gpt3.EditsResponse) {
	var responses []string
	for i := range resp.Choices {
		responses = append(responses, resp.Choices[i].Text, "\n\n###\n\n")
	}
	parameters.PromptCtx = responses
	log := strings.Join(responses, "")
	reg := strings.ReplaceAll(log, "[]", "\n")
	node.layout.promptOutput.SetText(reg)
	node.layout.infoOutput.SetText(fmt.Sprintf("\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
		resp.Created,
		resp.Object,
		resp.Usage.CompletionTokens,
		resp.Usage.PromptTokens,
		resp.Usage.TotalTokens,
		resp.Choices[0].Index))
}
