// Package service section
package service

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"caos/model"
	"caos/util"
	"encoding/json"

	"github.com/PullRequestInc/go-gpt3"
)

// EventManager - Log event service
type EventManager struct {
	pool model.PoolProperties
}

// ExportTraining - Export training in JSON format
func (c *EventManager) ExportTraining(session []model.TrainingSession) {
	raw, _ := json.MarshalIndent(session, "", "\u0009")
	out := util.ConstructTsPathFileTo("training", "json")
	out.WriteString(string(raw))
}

// saveLogSession - Save log session with actual detail
func (c *EventManager) saveLogSession() {
	if c.pool.Session != nil {
		raw, _ := json.MarshalIndent(c.pool.Session[len(c.pool.Session)-1], "", "\u0009")
		out := util.ConstructTsPathFileTo("log", "json")
		out.WriteString(string(raw))
	}
}

// clearSession - Clear all the pools
func (c *EventManager) clearSession() {
	c.pool.Event = nil
	c.pool.Session = nil
	c.pool.TrainingEvent = nil
	c.pool.TrainingSession = nil
}

// appendToSession - Add a set of events as a session
func (c *EventManager) appendToSession(id string, prompt model.HistoricalPrompt, train model.TrainingPrompt) {
	lEvent := model.HistoricalEvent{
		Timestamp: fmt.Sprint(time.Now().UnixMilli()),
		Event:     prompt,
	}

	c.pool.Event = append(c.pool.Event, lEvent)

	lSession := model.HistoricalSession{
		ID:      id,
		Session: []model.HistoricalEvent{lEvent},
	}

	c.pool.Session = append(c.pool.Session, lSession)

	event := model.TrainingEvent{
		Timestamp: fmt.Sprint(time.Now().UnixMilli()),
		Event:     train,
	}

	c.pool.TrainingEvent = append(c.pool.TrainingEvent, event)

	session := model.TrainingSession{
		ID:      id,
		Session: []model.TrainingEvent{event},
	}

	c.pool.TrainingSession = append(c.pool.TrainingSession, session)

	c.saveLogSession()
}

// appendToLayout - Append and visualize content in console page view
func (c *EventManager) appendToLayout(responses []string) {
	log := strings.Join(responses, "")
	node.layout.promptOutput.SetText(log)
}

// appendToChoice - Append choice to response
func (c *EventManager) appendToChoice(comp *gpt3.CompletionResponse, edit *gpt3.EditsResponse, search *gpt3.EmbeddingsResponse, chat *gpt3.ChatCompletionResponse, predict *model.Predict) []string {
	var responses []string
	responses = append(responses, "\n")
	if comp != nil && edit == nil && search == nil && chat == nil {
		if node.controller.currentAgent.preferences.IsPromptStreaming && comp.Choices != nil {
			responses = append(responses, comp.Choices[0].Text, "\n\n###\n\n")
		} else {
			for i := range comp.Choices {
				responses = append(responses, comp.Choices[i].Text, "\n\n###\n\n")
			}
		}
	} else if edit != nil && comp == nil && search == nil && chat == nil {
		for i := range edit.Choices {
			responses = append(responses, edit.Choices[i].Text, "\n\n###\n\n")
		}
	} else if chat != nil && comp == nil && search == nil && edit == nil {
		if node.controller.currentAgent.preferences.IsPromptStreaming && chat.Choices != nil {
			responses = append(responses, chat.Choices[0].Message.Content, "\n\n###\n\n")
		} else {
			for i := range chat.Choices {
				responses = append(responses, chat.Choices[i].Message.Content, "\n\n###\n\n")
			}
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

// appendToModel - Append conversation to session model
func (c *EventManager) appendToModel(transform model.TemplateProperties, header model.EngineProperties, body model.PromptProperties, predictBody model.PredictProperties, completion []string) (model.TrainingPrompt, model.HistoricalPrompt) {
	var modelTrainer model.TrainingPrompt
	var modelPrompt model.HistoricalPrompt

	modelTrainer = model.TrainingPrompt{
		Prompt:     body.Input,
		Completion: completion,
	}

	modelPrompt = model.HistoricalPrompt{
		Header:         header,
		Body:           body,
		PredictiveBody: predictBody,
		Template:       transform,
	}

	return modelTrainer, modelPrompt
}

// checkNewSession - Evaluate a new session
func (c *EventManager) checkNewSession() {
	if node.controller.currentAgent.preferences.IsNewSession {
		c.clearSession()
	}
}

// LogChatCompletion - Chat response details in a .json file
func (c *EventManager) LogChatCompletion(chain model.TemplateProperties, header model.EngineProperties, body model.PromptProperties, resp *gpt3.ChatCompletionResponse, cresp *gpt3.ChatCompletionStreamResponse) {
	c.checkNewSession()

	var modelTrainer model.TrainingPrompt
	var modelPrompt model.HistoricalPrompt

	if resp != nil && cresp == nil {
		for i := range resp.Choices {
			body.Content = []string{resp.Choices[i].Message.Content}
			modelTrainer, modelPrompt = c.appendToModel(chain, header, body, model.PredictProperties{}, []string{resp.Choices[i].Message.Content})
		}

		c.appendToSession(resp.ID, modelPrompt, modelTrainer)
		node.controller.currentAgent.preferences.CurrentID = resp.ID
	} else if cresp != nil && resp == nil {
		for i := range cresp.Choices {
			body.Content = []string{cresp.Choices[i].Delta.Content}
			modelTrainer, modelPrompt = c.appendToModel(chain, header, body, model.PredictProperties{}, []string{cresp.Choices[i].Delta.Content})
		}

		c.appendToSession(cresp.ID, modelPrompt, modelTrainer)
		node.controller.currentAgent.preferences.CurrentID = cresp.ID
	}
}

// LogGeneralCompletion - Response details in .json format
func (c *EventManager) LogGeneralCompletion(header model.EngineProperties, body model.PromptProperties, resp []string, id string) {
	c.checkNewSession()

	body.Content = resp
	modelTrainer, modelPrompt := c.appendToModel(model.TemplateProperties{}, header, body, model.PredictProperties{}, resp)

	c.appendToSession(id, modelPrompt, modelTrainer)
	node.controller.currentAgent.preferences.CurrentID = id
}

// LogPredict - ResponseDetails in a .json file
func (c *EventManager) LogPredict(header model.EngineProperties, body model.PredictProperties, resp *model.PredictResponse) {
	c.checkNewSession()

	var modelTrainer model.TrainingPrompt
	var modelPrompt model.HistoricalPrompt

	for i := range resp.Documents {
		PredictProperties := model.PredictProperties{
			Input:   body.Input,
			Details: *resp,
		}

		modelTrainer, modelPrompt = c.appendToModel(model.TemplateProperties{}, header, model.PromptProperties{}, PredictProperties, []string{fmt.Sprintf("%v", resp.Documents[i])})
	}

	c.appendToSession(node.controller.currentAgent.preferences.CurrentID, modelPrompt, modelTrainer)
}

// VisualLogCompletion - Chat response details
func (c *EventManager) VisualLogCompletion(resp *gpt3.CompletionResponse, cresp *gpt3.ChatCompletionResponse, sresp *gpt3.ChatCompletionStreamResponse) {
	if resp != nil && cresp == nil && sresp == nil {
		c.appendToLayout(c.appendToChoice(resp, nil, nil, nil, nil))

		for i := range resp.Choices {
			node.layout.infoOutput.SetText(
				fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nToken probs: %v \nToken top: %v\nFinish reason: %v\nIndex: %v\n",
					resp.ID,
					resp.Model,
					resp.Created,
					resp.Object,
					resp.Usage.CompletionTokens,
					resp.Usage.PromptTokens,
					resp.Usage.TotalTokens,
					resp.Choices[i].LogProbs.TokenLogprobs,
					resp.Choices[i].LogProbs.TopLogprobs,
					resp.Choices[i].FinishReason,
					resp.Choices[i].Index))
		}
	} else if cresp != nil && sresp == nil && resp == nil {
		c.appendToLayout(c.appendToChoice(nil, nil, nil, cresp, nil))

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
	} else if sresp != nil && cresp == nil && resp == nil {
		for i := range sresp.Choices {
			node.layout.infoOutput.SetText(
				fmt.Sprintf("ID: %v\nModel: %v\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nFinish reason: %v\nIndex: %v \n",
					sresp.ID,
					sresp.Model,
					sresp.Created,
					sresp.Object,
					sresp.Usage.CompletionTokens,
					sresp.Usage.PromptTokens,
					sresp.Usage.TotalTokens,
					sresp.Choices[i].FinishReason,
					sresp.Choices[i].Index))
		}
	}
}

// VisualLogEdit - Log edited response details
func (c *EventManager) VisualLogEdit(resp *gpt3.EditsResponse) {
	c.appendToLayout(c.appendToChoice(nil, resp, nil, nil, nil))
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
func (c *EventManager) VisualLogEmbedding(resp *gpt3.EmbeddingsResponse) {
	c.appendToLayout(c.appendToChoice(nil, nil, resp, nil, nil))
	for i := range resp.Data {
		node.layout.infoOutput.SetText(fmt.Sprintf("Object: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
			resp.Object,
			resp.Usage.PromptTokens,
			resp.Usage.TotalTokens,
			resp.Data[i].Index))
	}
}

// VisualLogPredict - Log predicted response details
func (c *EventManager) VisualLogPredict(resp *model.PredictResponse) {
	var buffer []string
	for i := range resp.Documents {
		c.appendToLayout(c.appendToChoice(nil, nil, nil, nil, &resp.Documents[i]))

		details := fmt.Sprintf("Average probability: %v\nCompletely generated probability: %v\nOverall burstiness: %v",
			resp.Documents[i].AverageProb,
			resp.Documents[i].CompletelyProb,
			resp.Documents[i].OverallBurstiness)
		buffer = append(buffer, details, "\n")

		for o := range resp.Documents[i].Paragraphs {
			paragraphs := fmt.Sprintf("\nCompletely generated probability: %v\nIndex: %v\nNumber of sentences: %v",
				resp.Documents[i].Paragraphs[o].CompletelyProb,
				resp.Documents[i].Paragraphs[o].Index,
				resp.Documents[i].Paragraphs[o].NumberSentences)
			buffer = append(buffer, paragraphs, "\n")
		}

		for o := range resp.Documents[i].Sentences {
			sentence := fmt.Sprintf("\nGenerated probability:%v\nPerplexity: %v\nSentence: %v",
				resp.Documents[i].Sentences[o].GeneratedProb,
				resp.Documents[i].Sentences[o].Perplexity,
				resp.Documents[i].Sentences[o].Sentence)
			buffer = append(buffer, sentence, "\n")
		}
	}

	inline := fmt.Sprintf("%v", buffer)
	node.layout.infoOutput.SetText(util.RemoveWrapper(inline))
}

// LogClient - Log client context
func (c *EventManager) LogClient(client Agent) {
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
	fmt.Printf("\n-------  /|_/|  ---------------------------")
	fmt.Printf("\n------- ( o.o ) ---------------------------")
	fmt.Printf("\n-------  > ^ <  ---------------------------")
	fmt.Printf("\n-------------------------------------------\n")
	fmt.Printf("ID name can be changed in PROFILE section\nmore information in: https://github.com/dabumana/caos\nClient ID: %v", client.id)
	fmt.Printf("\n-------------------------------------------\n")
	fmt.Print(`This software is provided "as is" and any expressed or implied warranties, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose are disclaimed. In no event shall the author or contributors be liable for any direct, indirect, incidental, special, exemplary, or consequential.`)
	fmt.Printf("\n-------------------------------------------\n")
}

// LogEngine - Log current engine
func (c *EventManager) LogEngine(client Agent) {
	node.layout.metadataOutput.SetText(
		fmt.Sprintf("Model: %v\nRole: %v\nTemperature: %v\nTopp: %v\nFrequency penalty: %v\nPresence penalty: %v\nPrompt: %v\nInstruction: %v\nProbabilities: %v\nResults: %v\nMax tokens: %v\n",
			client.EngineProperties.Model,
			client.EngineProperties.Role,
			client.EngineProperties.Temperature,
			client.EngineProperties.TopP,
			client.EngineProperties.FrequencyPenalty,
			client.EngineProperties.PresencePenalty,
			client.PromptProperties.Input,
			client.PromptProperties.Instruction,
			client.PromptProperties.Probabilities,
			client.PromptProperties.Results,
			client.preferences.MaxTokens))
}

// LogPredictEngine - Log current predict engine
func (c *EventManager) LogPredictEngine(client Agent) {
	var out string
	for i := range client.PredictProperties.Details.Documents {
		if client.PredictProperties.Details.Documents[i].AverageProb >= 0.5 {
			out = "Probably generated by AI"
		} else {
			out = "Mostly human generated content"
		}

		node.layout.metadataOutput.SetText(
			fmt.Sprintf("Model: %v\nAverage Prob: %v\nCompletely Prob: %v\noversall burstiness: %v\n---\n%v\n",
				client.EngineProperties.Model,
				client.PredictProperties.Details.Documents[i].AverageProb,
				client.PredictProperties.Details.Documents[i].CompletelyProb,
				client.PredictProperties.Details.Documents[i].OverallBurstiness,
				out))
	}
}

// Errata - Generic error method
func (c *EventManager) Errata(err error) {
	if flag.Lookup("test.v") == nil {
		if err != nil {
			node.layout.infoOutput.SetText(err.Error())
			node.layout.promptArea.SetPlaceholder("An error was found or the response was not complete, just press CTRL+SPACE or CMD+SPACE to repeat it.")
		} else {
			node.layout.promptArea.SetPlaceholder("Type here...")
		}

		node.layout.promptArea.SetText("", true)
	}
}
