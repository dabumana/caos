// Package service section
package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"caos/model"
	"caos/util"

	"github.com/PullRequestInc/go-gpt3"
)

// EventManager - Log event service
type EventManager struct {
	event   model.HistoricalEvent
	session model.HistoricalSession
	pool    model.PoolProperties
}

// SaveTraining - Export training in JSON format
func (c EventManager) SaveTraining() {
	raw, _ := json.MarshalIndent(c.pool.TrainingSession, "", "\u0009")
	out := util.ConstructPathFileTo("training", "json")
	out.WriteString(string(raw))
}

// SaveLog - Save log with actual historic detail
func (c EventManager) SaveLog() {
	if c.pool.Session != nil {
		raw, _ := json.MarshalIndent(c.pool.Session[len(c.pool.Session)-1], "", "\u0009")
		out := util.ConstructPathFileTo("log", "json")
		out.WriteString(string(raw))
	}
}

// ClearSession - Clear all the pools
func (c EventManager) ClearSession() {
	c.pool.Event = nil
	c.pool.Session = nil
	c.pool.TrainingEvent = nil
	c.pool.TrainingSession = nil
}

// AppendToSession - Add a set of events as a session
func (c EventManager) AppendToSession(id string, prompt model.HistoricalPrompt, train model.TrainingPrompt) {
	historical := model.HistoricalEvent{
		Timestamp: fmt.Sprint(time.Now().UnixMilli()),
		Event:     prompt,
	}

	c.pool.Event = append(c.pool.Event, historical)

	lsession := model.HistoricalSession{
		ID:      id,
		Session: []model.HistoricalEvent{historical},
	}

	c.pool.Session = append(c.pool.Session, lsession)

	if node.controller.currentAgent.preferences.IsTraining {

		event := model.HistoricalTrainingEvent{
			Timestamp: c.event.Timestamp,
			Event:     train,
		}

		c.pool.TrainingEvent = append(c.pool.TrainingEvent, event)

		session := model.HistoricalTrainingSession{
			ID:      c.session.ID,
			Session: []model.HistoricalTrainingEvent{event},
		}

		c.pool.TrainingSession = append(c.pool.TrainingSession, session)
	}

	c.SaveLog()
}

// AppendToLayout - Append and visualize content in console page view
func (c EventManager) AppendToLayout(responses []string) {
	log := strings.Join(responses, "")
	node.layout.promptOutput.SetText(log)
}

// AppendToChoice - Append choice to response
func (c EventManager) AppendToChoice(comp *gpt3.CompletionResponse, edit *gpt3.EditsResponse, search *gpt3.EmbeddingsResponse, chat *gpt3.ChatCompletionResponse, predict *model.Predict) []string {
	var responses []string
	responses = append(responses, "\n")
	if comp != nil && edit == nil && search == nil && chat == nil {
		for i := range comp.Choices {
			responses = append(responses, comp.Choices[i].Text, "\n\n###\n\n")
		}
	} else if edit != nil && comp == nil && search == nil && chat == nil {
		for i := range edit.Choices {
			responses = append(responses, edit.Choices[i].Text, "\n\n###\n\n")
		}
	} else if chat != nil && comp == nil && search == nil && edit == nil {
		for i := range chat.Choices {
			responses = append(responses, chat.Choices[i].Message.Content, "\n\n###\n\n")
		}
	} else if predict != nil && edit == nil && comp == nil && search == nil {
		for i := range predict.Sentences {
			responses = append(responses, predict.Sentences[i].Sentence, "\n\n###\n\n")
		}
	} else {
		for i := range search.Data {
			responses = append(responses, fmt.Sprintf("%v", search.Data[i]), "\n\n###\n\n")
		}
	}
	return responses
}

// LogChatCompletion - Chat response details in a .json file
func (c EventManager) LogChatCompletion(header model.EngineProperties, body model.PromptProperties, resp *gpt3.ChatCompletionResponse, cresp *gpt3.ChatCompletionStreamResponse) {
	if node.controller.currentAgent.preferences.IsNewSession {
		c.ClearSession()
		node.controller.currentAgent.preferences.IsNewSession = false
	}

	if resp != nil && cresp == nil {
		for i := range resp.Choices {
			body.Content = []string{resp.Choices[i].Message.Content}

			modelTrainer := model.TrainingPrompt{
				Prompt:     body.PromptContext,
				Completion: []string{resp.Choices[i].Message.Content},
			}

			modelPrompt := model.HistoricalPrompt{
				Header: header,
				Body:   body,
			}

			c.AppendToSession(resp.ID, modelPrompt, modelTrainer)

			node.controller.currentAgent.preferences.CurrentID = resp.ID
		}
	} else if cresp != nil && resp == nil {
		for i := range cresp.Choices {
			body.Content = []string{cresp.Choices[i].Delta.Content}

			modelTrainer := model.TrainingPrompt{
				Prompt:     body.PromptContext,
				Completion: []string{cresp.Choices[i].Delta.Content},
			}

			modelPrompt := model.HistoricalPrompt{
				Header: header,
				Body:   body,
			}

			c.AppendToSession(cresp.ID, modelPrompt, modelTrainer)

			node.controller.currentAgent.preferences.CurrentID = cresp.ID
		}
	}
}

// LogCompletion - Response details in a .json file
func (c EventManager) LogCompletion(header model.EngineProperties, body model.PromptProperties, resp *gpt3.CompletionResponse) {
	if node.controller.currentAgent.preferences.IsNewSession {
		c.ClearSession()
		node.controller.currentAgent.preferences.IsNewSession = false
	}

	for i := range resp.Choices {
		body.Content = []string{resp.Choices[i].Text}

		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{resp.Choices[i].Text},
		}

		modelPrompt := model.HistoricalPrompt{
			Header: header,
			Body:   body,
		}

		c.AppendToSession(resp.ID, modelPrompt, modelTrainer)
		node.controller.currentAgent.preferences.CurrentID = resp.ID
	}
}

// LogEdit - Response details in a .json file
func (c EventManager) LogEdit(header model.EngineProperties, body model.PromptProperties, resp *gpt3.EditsResponse) {

	for i := range resp.Choices {
		body.Content = []string{resp.Choices[i].Text}

		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{resp.Choices[i].Text},
		}

		modelPrompt := model.HistoricalPrompt{
			Header: header,
			Body:   body,
		}

		c.AppendToSession(node.controller.currentAgent.preferences.CurrentID, modelPrompt, modelTrainer)
	}
}

// LogEmbedding - Response details in a .json file
func (c EventManager) LogEmbedding(header model.EngineProperties, body model.PromptProperties, resp *gpt3.EmbeddingsResponse) {
	for i := range resp.Data {
		body.Content = []string{resp.Data[i].Object}

		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{fmt.Sprintf("%v", resp.Data[i])},
		}

		modelPrompt := model.HistoricalPrompt{
			Header: header,
			Body:   body,
		}

		c.AppendToSession(node.controller.currentAgent.preferences.CurrentID, modelPrompt, modelTrainer)
	}
}

// LogPredict - ResponseDetails in a .json file
func (c EventManager) LogPredict(header model.EngineProperties, body model.PromptProperties, predict *model.PredictProperties, resp *model.PredictResponse) {
	for i := range resp.Documents {
		predict.Details.Documents = append(predict.Details.Documents, resp.Documents[i])
		modelTrainer := model.TrainingPrompt{
			Prompt:     predict.Input,
			Completion: []string{fmt.Sprintf("%v", resp.Documents[i])},
		}

		modelPrompt := model.HistoricalPrompt{
			Header: header,
			Body:   body,
		}

		c.AppendToSession(node.controller.currentAgent.preferences.CurrentID, modelPrompt, modelTrainer)
	}
}

// VisualLogChatCompletion - Chat response details
func (c EventManager) VisualLogChatCompletion(resp *gpt3.ChatCompletionResponse, cresp *gpt3.ChatCompletionStreamResponse) {
	if resp != nil && cresp == nil {
		if !node.controller.currentAgent.preferences.IsPromptStreaming {
			c.AppendToLayout(c.AppendToChoice(nil, nil, nil, resp, nil))
		}
		for i := range resp.Choices {
			node.layout.infoOutput.SetText(
				fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nIndex: %v \n",
					resp.ID,
					resp.Model,
					resp.Created,
					resp.Object,
					resp.Usage.CompletionTokens,
					resp.Usage.PromptTokens,
					resp.Usage.TotalTokens,
					resp.Choices[i].FinishReason,
					resp.Choices[i].Index))
		}
	} else if cresp != nil && resp == nil {
		for i := range cresp.Choices {
			node.layout.infoOutput.SetText(
				fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nIndex: %v \n",
					cresp.ID,
					cresp.Model,
					cresp.Created,
					cresp.Object,
					cresp.Usage.CompletionTokens,
					cresp.Usage.PromptTokens,
					cresp.Usage.TotalTokens,
					cresp.Choices[i].FinishReason,
					cresp.Choices[i].Index))
		}
	}
}

// VisualLogCompletion - Response details
func (c EventManager) VisualLogCompletion(resp *gpt3.CompletionResponse) {
	if !node.controller.currentAgent.preferences.IsPromptStreaming {
		c.AppendToLayout(c.AppendToChoice(resp, nil, nil, nil, nil))
	}
	for i := range resp.Choices {
		node.layout.infoOutput.SetText(
			fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nToken probs: %v \nToken top: %v\nIndex: %v\n",
				resp.ID,
				resp.Model,
				resp.Created,
				resp.Object,
				resp.Usage.CompletionTokens,
				resp.Usage.PromptTokens,
				resp.Usage.TotalTokens,
				resp.Choices[i].FinishReason,
				resp.Choices[i].LogProbs.TokenLogprobs,
				resp.Choices[i].LogProbs.TopLogprobs,
				resp.Choices[i].Index))
	}
}

// VisualLogEdit - Log edited response details
func (c EventManager) VisualLogEdit(resp *gpt3.EditsResponse) {
	if !node.controller.currentAgent.preferences.IsPromptStreaming {
		c.AppendToLayout(c.AppendToChoice(nil, resp, nil, nil, nil))
	}
	for i := range resp.Choices {
		node.layout.infoOutput.SetText(fmt.Sprintf("Created: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
			resp.Created,
			resp.Object,
			resp.Usage.CompletionTokens,
			resp.Usage.PromptTokens,
			resp.Usage.TotalTokens,
			resp.Choices[i].Index))
	}
}

// VisualLogEmbedding - Log embedding response details
func (c EventManager) VisualLogEmbedding(resp *gpt3.EmbeddingsResponse) {
	if !node.controller.currentAgent.preferences.IsPromptStreaming {
		c.AppendToLayout(c.AppendToChoice(nil, nil, resp, nil, nil))
	}
	for i := range resp.Data {
		node.layout.infoOutput.SetText(fmt.Sprintf("Object: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
			resp.Object,
			resp.Usage.PromptTokens,
			resp.Usage.TotalTokens,
			resp.Data[i].Index))
	}
}

// VisualLogPredict - Log predicted response details
func (c EventManager) VisualLogPredict(resp *model.PredictResponse) {
	var buffer []string
	for i := range resp.Documents {
		c.AppendToLayout(c.AppendToChoice(nil, nil, nil, nil, &resp.Documents[i]))

		details := fmt.Sprintf("Average probability: %v\nCompletely generated probability: %v\nOverall burstiness: %v",
			resp.Documents[i].AverageProb,
			resp.Documents[i].CompletelyProb,
			resp.Documents[i].OverallBurstiness)
		buffer = append(buffer, details, "\n")

		for o := range resp.Documents[i].Paragraphs {
			paragraphs := fmt.Sprintf("Completely generated probability: %v\nIndex: %v\nNumber of sentences: %v",
				resp.Documents[i].Paragraphs[o].CompletelyProb,
				resp.Documents[i].Paragraphs[o].Index,
				resp.Documents[i].Paragraphs[o].NumberSentences)
			buffer = append(buffer, paragraphs, "\n")
		}

		for o := range resp.Documents[i].Sentences {
			sentence := fmt.Sprintf("Generated probability:%v\nPerplexity: %v\nSentence: %v",
				resp.Documents[i].Sentences[o].GeneratedProb,
				resp.Documents[i].Sentences[o].Perplexity,
				resp.Documents[i].Sentences[o].Sentence)
			buffer = append(buffer, sentence, "\n")
		}
	}
	inline := fmt.Sprintf("%v", buffer)
	output := strings.ReplaceAll(inline, "[", "")
	output = strings.ReplaceAll(output, "]", "")
	node.layout.infoOutput.SetText(output)
}

// LogClient - Log client context
func (c EventManager) LogClient(client Agent) {
	fmt.Print(`
 _______________________________________________________________________________________________________________________________________________________________________________________________________
|............................................................................................C.A.O.S....................................................................................................|	
|_______________________________________________________________________________________________________________________________________________________________________________________________________|
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&&&&&&&&&&&&&&&&&&&&&&&                   &&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&& *&&&&&&&&&&                        /&&               &&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&/                                                  #&&           %&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&                                                                &&         &&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&                                                                           &.      &&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&          &                     .................................                  ,&     &&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&                             ..............................................                &    &&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&                                ................................................&.....            &  &&&&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&                                .. ...................................................&&.......          &    &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&                     &              ...&&...............&...................................&........         &       &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&                    &        .     ..&..  ......&.........&......................................&..........       &..  & %&&&&&&&&&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@&          ..        &   ............&...    ....&...........&.........................&...............&..........       &&                     &&&@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@&        ...#        &   ............&...........&..   .........&......%...................&................&...........     /#...                         &&@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@&      ....&       .&   ...........,.............#..   ...   ...,&.......&....................&.................&...........     &.....                            &&@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@&     ...%&&&     ..&   ............&............&...............&..&.......&.....................&...........&.....&............    &............                         &&@@@@@@@
@@@@@@@@@@@@@@@@@@&    ..&&@@&     ..#.   ...........&............&..............%&&....&..  ...&.........&............/.................%............    *&..............                         &&@@@
@@@@@@@@@@@@@@@&*  ..&@@@@@&    ...&    ...........&...........&&.............&&&&......&.......&... ...................................................    #...............                          &&
@@@@@@@@@@@@@&  .&@@@@@@@&   ....&    ...........&..........,&&...........&*..&.........&%......&... .&./..&............&.............,.....#............   &..................                        *
@@@@@@@@@@@& &@@@@@@@@@&  ...&&%    ..........&&(.....&&&&&.&........&&....*(............&......,&.......................&.............%.....&............. &&...................                       
@@@@@@@@@@&@@@@@@@@@@& ..#&@&#    .#........&&&.......&&..&...*&&&.......&...............&.......&..........&...................,.......%.....&........... &  &....................                     
@@@@@@@@@@@@@@@@@@@&..&&@@&&    .&........&&&......&,...&...&%.....&..&..................&.......&.........../............&.....&............../......... ...  &..&.................                   &
@@@@@@@@@@@@@@@@@&.&@@@@@&     .&.......&&.&...&&&&%..&&&..........&......................&.....,.&..........&.............&...................%.............. &....&.................                &@
@@@@@@@@@@@@@@@@&@@@@@@&&    ..#.....,&&..&&&.............#.....%.........................&.....&..&........../............&............%.#.....&.............  &....&.................              &@@
@@@@@@@@@@@@@@@@@@@@@@&    .........,&&&&&&&&&&&&&&&.........#.............................&....&..&*.........&.............&..............&.....&............  &......&................            &@@@
@@@@@@@@@@@@@@@@@@@@@&    ..#......&&..(&.......(&&&&&&&&................................../....&...&..........&............&...............&.....,...........  &.......&................         &@@@@@
@@@@@@@@@@@@@@@@@@@@&    ..&.....&&&&..  &     &&     *&&&&.................................&...&..#(&&........&..................................&...........  &.......................        &@@@@@@@
@@@@@@@@@@@@@@@@@@@     ..%.....&&&&&.   ,&&   &       &  &&&.&..............................&..&....../.,&.....&............&.....&.........&.....&.........  (*........&.............       &@@@@@@@@@
@@@@@@@@@@@@@@@@@&     .......&&&,.&&      .......&   &&     &&.&.............................,.&.......&.....&.&&...........&.....................&......... #..*..................       &&@@@@@@@@@@@
@@@@@@@@@@@@@@@@&    ........&&&.&.&.     ..............,&   &#&..............................&.&.&......&......&&...........*................&.....&......  ...&...............        &&@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@    ....&...&&&.&*%&  &   .....................................................&&.......&.&.....&.&.................................&......&....&.........             #&@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@    ....&...&&&....&.   & .......,&....................................,&#....%&&&..........&....&.&.................&..........&....&....&......&.........         %   &@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@.  .....&...@&......&    %,..............&..............................*&&&&&&&&&&&&&#.*/....&*&....&..........,.....&..........&.....&....&....&..........            &@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@&  ....&&&..@@......&      ....,&...................  .............................&&&&&&&&&&&&..&.*..&..........*.....&..........#.....&........&&..........           &@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@&  ....&@@..@@......&      ........................&%    ......................    &        &&&&&&&&.,,.%........,/.....&................&.......(&...........          &@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@ ....&@@@&.&@&......&      ......................  .................................&,     &&   &&&&&&%.&........&......&...........%....&......*&&...........        &@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@ ...&@@@@@.&@&......&      ...............................................................&&&       &&&&&&&.......&......&...........&.....#....%.&............      &@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@&..&@@@@@@&&@@&......&      ....................................................................&    &...&&&&.&&...&......%...........&&....&.....&............   &&&@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@./@@@@@@@@&@@&.......*     ........................................................(&...............&..,&&&&&&&&..&&..............#...&.&...&.*..&.&&........ /&& &&@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@&@@@@@@@@@@@@@/.......&     ............................................................(&&................&.......&&..............&...&.&...(...&..&..&...&&&..   &@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@&........&     ...........&..........................................................,#,.....&..&....&&#.....&........&..%...&..#.........&......&.   &@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@&......&@&      ............&...............................................................&%..&....&.......&........&..&....*.&&....&...&......#.   &@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@.....&&@@@&     ................&..........................................................%&...&...&...&...&.........&..&....&.&.........&........   @@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@&....&@@@@@@&    .........................................................................#&&.......&....&...&.................&((....&...&......&.   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@&..&@@@@@@@@@&    ............................&,...........................................&.......&....&&..(.........&..&...&%.&*.../....%........   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@*.&@@@@@@@@@@@&    .....................................,&&*...%&&........................&......#&....&.&..&.........&..&.&.&..&........&......&.   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@.&@@@@@@@@@@@@@&    .....................................................................&......./...&(..&.&%#........&.&&..&...........&.........   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@&@@@@@@@@@@@@@@@@&    ..................................................................&...........&....%...&.......%..&..&...................&.   (@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&    ................................................................&..........&.......&..........&.*..*......................   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&    .............................................................&.........&.............&......&.&........................   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&   ...........................................................&........&.&#.....&....&.&.....&&/............,.........*.    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@    ........................................................#......&/............&&@@& &....&&&...........&............   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@    ..................................................&%&......&........&&&&@@@@@@@%  ...&.&........................   &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%   .................................................&....&  &/ @@@@@@@@@@@@@@@@&  .&....&.....*...&............. &  @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&      .............................&&%..............&@@@&  &@& @@@@@@@@@@@@&   .&.&..........&........... ...&  &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&&&&&&&&&&&&&&&&&&&................&..&.&@@@@@&  @@@@@@@@@@@@@@&   ..............#........... %.. & &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&&&&&&&&&&...................&.%...&@@@@@@@& &@@@@@@@@@@@@&   ....&...,....&...........   ...& .@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@  .................................&@@@@@@@@@@@@@@@@@@@@@&   .........................   ....& &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@  ................................#&&&&@@@@&@@@@@@@@@@@@@    ......&....&............   &....#&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&  ................................&((((((((((@@@@@@@@@@@    .................(......   &&...& &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@& ......................(((((((((((((#&&%(((((@@@@@@@@@@    .....&....&.............   &@&...&&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@  ............&%(((((((((((((((((((((((((((((@@@@@@@@@@& & ..........&.............   &@@@&..&&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&  ......&(((((((((((((((((((((((((((((((((((&@@@@@@@@@@/&  ...&.....,.....&......    &@@@@&...&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&.  .&((((((((((((((((((((((((#(((((((((((((((((&@@@@@@@@& &  ..&&.../......&......    &@@@@@@&..&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&%(((((((((((((((((((((((((((((((((((((((((((((((((&@@@@@@@&   ..&/(..&......&......   *@@@@@@@@@@..&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&((((((((((((     ......(((((((((((((((((((((((((((((((((&@@@@@@&&  .& &..&*...../.....    &@@@@@@@@@@@@&.&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&((((((((((((%&& .....&&&#(((.,..(.(((((((((((((((((((((((((((((&@@@@@&  .  &..&......&....    &@@@@@@@@@@@@@@@@&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&&((((((&&     ..............&&&((((.,./(((((((((((((((((((((((((((((&@@@@& &  (.&......*....    @@@@@@@@@@@@@@@@@@@&&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
	`)
	fmt.Printf("\n-------------------------------------------\n")
	fmt.Printf("Context: %v\nClient: %v\n", client.ctx, client.client)
	fmt.Printf("-------------------------------------------\n")
	fmt.Print(`This software is provided "as is" and any expressed or implied warranties, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose are disclaimed. In no event shall the author or contributors be liable for any direct, indirect, incidental, special, exemplary, or consequential.`)
	fmt.Printf("\n-------------------------------------------\n")
}

// LogEngine - Log current engine
func (c EventManager) LogEngine(client Agent) {
	node.layout.metadataOutput.SetText(
		fmt.Sprintf("Model: %v\nRole: %v\nTemperature: %v\nTopp: %v\nFrequency penalty: %v\nPresence penalty: %v\nPrompt: %v\nInstruction: %v\nProbabilities: %v\nResults: %v\nMax tokens: %v\n",
			client.engineProperties.Model,
			client.engineProperties.Role,
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

// LogPredictEngine - Log current predict engine
func (c EventManager) LogPredictEngine(client Agent) {
	var out string
	for i := range client.predictProperties.Details.Documents {
		if client.predictProperties.Details.Documents[i].AverageProb >= 0.5 {
			out = "Probably generated by AI"
		} else {
			out = "Mostly human generated content"
		}

		node.layout.metadataOutput.SetText(
			fmt.Sprintf("Model: %v\nAverage Prob: %v\nCompletely Prob: %v\noversall burstiness: %v\n---\n%v\n",
				client.engineProperties.Model,
				client.predictProperties.Details.Documents[i].AverageProb,
				client.predictProperties.Details.Documents[i].CompletelyProb,
				client.predictProperties.Details.Documents[i].OverallBurstiness,
				out))
	}
}

// Errata - Generic error method
func (c EventManager) Errata(err error) {
	if err != nil {
		node.controller.currentAgent.preferences.IsNewSession = true
		node.layout.infoOutput.SetText(err.Error())
		node.layout.promptArea.SetPlaceholder("An error was found repeat your request and press CTRL+SPACE or CMD+SPACE")
	} else {
		node.layout.promptArea.SetPlaceholder("Type here...")
	}

	node.controller.currentAgent.preferences.IsLoading = false
	node.layout.promptArea.SetText("", true)
}
