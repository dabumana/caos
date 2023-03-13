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
	event     model.HistoricalEvent
	session   model.HistoricalSession
	pool      model.PoolProperties
	isRunning bool
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
	log = strings.ReplaceAll(log, "[]", "\n")
	node.layout.promptOutput.SetText(log)
}

// AppendToChoice - Append choice to response
func (c EventManager) AppendToChoice(comp *gpt3.CompletionResponse, edit *gpt3.EditsResponse, search *gpt3.EmbeddingsResponse, chat *gpt3.ChatCompletionResponse, schat *gpt3.ChatCompletionStreamResponse) []string {
	var responses []string
	responses = append(responses, "\n")
	if comp != nil && edit == nil && search == nil && chat == nil && schat == nil {
		for i := range comp.Choices {
			responses = append(responses, comp.Choices[i].Text, "\n\n###\n\n")
		}
	} else if edit != nil && comp == nil && search == nil && chat == nil && schat == nil {
		for i := range edit.Choices {
			responses = append(responses, edit.Choices[i].Text, "\n\n###\n\n")
		}
	} else if chat != nil && comp == nil && search == nil && edit == nil && schat == nil {
		for i := range chat.Choices {
			responses = append(responses, chat.Choices[i].Message.Content, "\n\n###\n\n")
		}
	} else if schat != nil && comp == nil && search == nil && edit == nil && chat == nil {
		for i := range schat.Choices {
			responses = append(responses, schat.Choices[i].Delta.Content, "\n\n###\n\n")
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
		body.Content = []string{resp.Choices[0].Message.Content}

		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{resp.Choices[0].Message.Content},
		}

		modelPrompt := model.HistoricalPrompt{
			Header: header,
			Body:   body,
		}

		c.AppendToSession(resp.ID, modelPrompt, modelTrainer)

		node.controller.currentAgent.preferences.CurrentID = resp.ID

	} else if cresp != nil && resp == nil {
		body.Content = []string{cresp.Choices[0].Delta.Content}

		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{cresp.Choices[0].Delta.Content},
		}

		modelPrompt := model.HistoricalPrompt{
			Header: header,
			Body:   body,
		}

		c.AppendToSession(cresp.ID, modelPrompt, modelTrainer)

		node.controller.currentAgent.preferences.CurrentID = cresp.ID
	}
}

// LogCompletion - Response details in a .json file
func (c EventManager) LogCompletion(header model.EngineProperties, body model.PromptProperties, resp *gpt3.CompletionResponse) {
	if node.controller.currentAgent.preferences.IsNewSession {
		c.ClearSession()
		node.controller.currentAgent.preferences.IsNewSession = false
	}

	body.Content = []string{resp.Choices[0].Text}

	modelTrainer := model.TrainingPrompt{
		Prompt:     body.PromptContext,
		Completion: []string{resp.Choices[0].Text},
	}

	modelPrompt := model.HistoricalPrompt{
		Header: header,
		Body:   body,
	}

	c.AppendToSession(resp.ID, modelPrompt, modelTrainer)
	node.controller.currentAgent.preferences.CurrentID = resp.ID
}

// LogEdit - Response details in a .json file
func (c EventManager) LogEdit(header model.EngineProperties, body model.PromptProperties, resp *gpt3.EditsResponse) {
	body.Content = []string{resp.Choices[0].Text}

	modelTrainer := model.TrainingPrompt{
		Prompt:     body.PromptContext,
		Completion: []string{resp.Choices[0].Text},
	}

	modelPrompt := model.HistoricalPrompt{
		Header: header,
		Body:   body,
	}

	c.AppendToSession(node.controller.currentAgent.preferences.CurrentID, modelPrompt, modelTrainer)
}

// LogEmbedding - Response details in a .json file
func (c EventManager) LogEmbedding(header model.EngineProperties, body model.PromptProperties, resp *gpt3.EmbeddingsResponse) {
	body.Content = []string{resp.Data[0].Object}

	modelTrainer := model.TrainingPrompt{
		Prompt:     body.PromptContext,
		Completion: []string{fmt.Sprintf("%v", resp.Data[0])},
	}

	modelPrompt := model.HistoricalPrompt{
		Header: header,
		Body:   body,
	}

	c.AppendToSession(node.controller.currentAgent.preferences.CurrentID, modelPrompt, modelTrainer)
}

// VisualLogChatCompletion - Chat response details
func (c EventManager) VisualLogChatCompletion(resp *gpt3.ChatCompletionResponse, cresp *gpt3.ChatCompletionStreamResponse) {
	if resp != nil && cresp == nil {
		c.AppendToLayout(c.AppendToChoice(nil, nil, nil, resp, nil))
		node.layout.infoOutput.SetText(
			fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nIndex: %v \n",
				resp.ID,
				resp.Model,
				resp.Created,
				resp.Object,
				resp.Usage.CompletionTokens,
				resp.Usage.PromptTokens,
				resp.Usage.TotalTokens,
				resp.Choices[0].FinishReason,
				resp.Choices[0].Index))
	} else if cresp != nil && resp == nil {
		c.AppendToLayout(c.AppendToChoice(nil, nil, nil, nil, cresp))
		node.layout.infoOutput.SetText(
			fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nIndex: %v \n",
				cresp.ID,
				cresp.Model,
				cresp.Created,
				cresp.Object,
				cresp.Usage.CompletionTokens,
				cresp.Usage.PromptTokens,
				cresp.Usage.TotalTokens,
				cresp.Choices[0].FinishReason,
				cresp.Choices[0].Index))
	}
}

// VisualLogCompletion - Response details
func (c EventManager) VisualLogCompletion(resp *gpt3.CompletionResponse) {
	if !node.controller.currentAgent.preferences.IsPromptStreaming {
		c.AppendToLayout(c.AppendToChoice(resp, nil, nil, nil, nil))
	}

	node.layout.infoOutput.SetText(
		fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nToken probs: %v \nToken top: %v\nIndex: %v\n",
			resp.ID,
			resp.Model,
			resp.Created,
			resp.Object,
			resp.Usage.CompletionTokens,
			resp.Usage.PromptTokens,
			resp.Usage.TotalTokens,
			resp.Choices[0].FinishReason,
			resp.Choices[0].LogProbs.TokenLogprobs,
			resp.Choices[0].LogProbs.TopLogprobs,
			resp.Choices[0].Index))
}

// VisualLogEdit - Log edited response details
func (c EventManager) VisualLogEdit(resp *gpt3.EditsResponse) {
	c.AppendToLayout(c.AppendToChoice(nil, resp, nil, nil, nil))
	node.layout.infoOutput.SetText(fmt.Sprintf("Created: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
		resp.Created,
		resp.Object,
		resp.Usage.CompletionTokens,
		resp.Usage.PromptTokens,
		resp.Usage.TotalTokens,
		resp.Choices[0].Index))
}

// VisualLogEmbedding - Log embedding response details
func (c EventManager) VisualLogEmbedding(resp *gpt3.EmbeddingsResponse) {
	c.AppendToLayout(c.AppendToChoice(nil, nil, resp, nil, nil))
	node.layout.infoOutput.SetText(fmt.Sprintf("Object: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
		resp.Object,
		resp.Usage.PromptTokens,
		resp.Usage.TotalTokens,
		resp.Data[0].Index))
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

// LoaderStreaming - Generic loading animation
func (c EventManager) LoaderStreaming(in string) {
	go func() {
		fmt.Println(in + "/ \\ _ / \\ \n (  o . o  )")
	}()
}
