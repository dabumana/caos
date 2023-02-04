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

// EventManager - Log event service
type EventManager struct {
	event   model.HistoricalEvent
	session model.HistoricalSession
}

// SaveLog - Save log with actual historic detail
func (c EventManager) SaveLog() {
	if parameters.SessionPool != nil {
		raw, _ := json.MarshalIndent(parameters.SessionPool[len(parameters.SessionPool)-1], "", "\u0009")
		out := util.ConstructPathFileTo("log", "json")
		out.WriteString(string(raw))
	}
}

// ClearSession - Clear all the pools
func (c EventManager) ClearSession() {
	parameters.EventPool = nil
	parameters.SessionPool = nil
	parameters.TrainingEventPool = nil
	parameters.TrainingSessionPool = nil
}

// AppendToSession - Add a set of events as a session
func (c EventManager) AppendToSession(header *model.EngineProperties, body *model.PromptProperties, predict *model.PredictProperties, id string, train model.TrainingPrompt) {
	valid := func(eventType any) bool {
		return eventType != nil
	}

	if valid(header) {
		c.event.Event.Header = *header
	} else {
		c.event.Event.Header = *new(model.EngineProperties)
	}

	if valid(body) {
		c.event.Event.Body = *body
	} else {
		c.event.Event.Body = *new(model.PromptProperties)
	}

	if valid(predict) {
		c.event.Event.Predict = *predict
	} else {
		c.event.Event.Predict = *new(model.PredictProperties)
	}

	c.event.Timestamp = fmt.Sprint(time.Now().UnixMilli())

	parameters.EventPool = append(parameters.EventPool, c.event)

	c.session.ID = id
	c.session.Session = parameters.EventPool

	if parameters.IsTraining {

		event := model.HistoricalTrainingEvent{
			Timestamp: c.event.Timestamp,
			Event:     train,
		}

		parameters.TrainingEventPool = append(parameters.TrainingEventPool, event)

		session := model.HistoricalTrainingSession{
			ID:      c.session.ID,
			Session: []model.HistoricalTrainingEvent{event},
		}

		parameters.TrainingSessionPool = append(parameters.TrainingSessionPool, session)
	}

	parameters.SessionPool = append(parameters.SessionPool, c.session)

	c.SaveLog()
}

// AppendToLayout - Append and visualize content in console page view
func (c EventManager) AppendToLayout(responses []string) {
	parameters.PromptCtx = responses
	log := strings.Join(responses, "")
	log = strings.ReplaceAll(log, "[]", "\n")
	node.layout.promptOutput.SetText(log)
}

// AppendToChoice - Append choice to response
func (c EventManager) AppendToChoice(comp *gpt3.CompletionResponse, edit *gpt3.EditsResponse, search *gpt3.EmbeddingsResponse, predict *model.Predict) []string {
	var responses []string
	responses = append(responses, "\n")
	if comp != nil && edit == nil && search == nil && predict == nil {
		for i := range comp.Choices {
			responses = append(responses, comp.Choices[i].Text, "\n\n###\n\n")
		}
	} else if edit != nil && comp == nil && search == nil && predict == nil {
		for i := range edit.Choices {
			responses = append(responses, edit.Choices[i].Text, "\n\n###\n\n")
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

// LogCompletion - Response details in a .json file
func (c EventManager) LogCompletion(header *model.EngineProperties, body *model.PromptProperties, resp *gpt3.CompletionResponse) {
	if parameters.IsNewSession {
		c.ClearSession()
		parameters.IsNewSession = false
	}

	for i := range resp.Choices {
		predict := new(model.PredictProperties)
		modelTrainer := model.TrainingPrompt{
			Prompt:     body.PromptContext,
			Completion: []string{resp.Choices[i].Text},
		}

		c.AppendToSession(header, body, predict, resp.ID, modelTrainer)
	}

	parameters.CurrentID = resp.ID
}

// LogEdit - Response details in a .json file
func (c EventManager) LogEdit(header *model.EngineProperties, body *model.PromptProperties, resp *gpt3.EditsResponse) {
	predict := new(model.PredictProperties)
	modelTrainer := model.TrainingPrompt{
		Prompt:     body.PromptContext,
		Completion: []string{resp.Choices[0].Text},
	}

	c.AppendToSession(header, body, predict, parameters.CurrentID, modelTrainer)
}

// LogEmbedding - Response details in a .json file
func (c EventManager) LogEmbedding(header *model.EngineProperties, body *model.PromptProperties, resp *gpt3.EmbeddingsResponse) {
	predict := new(model.PredictProperties)
	modelTrainer := model.TrainingPrompt{
		Prompt:     body.PromptContext,
		Completion: []string{fmt.Sprintf("%v", resp.Data[0])},
	}

	c.AppendToSession(header, body, predict, parameters.CurrentID, modelTrainer)
}

// LogPredict - ResponseDetails in a .json file
func (c EventManager) LogPredict(predict *model.PredictProperties, resp *model.PredictResponse) {
	header := new(model.EngineProperties)
	body := new(model.PromptProperties)
	modelTrainer := model.TrainingPrompt{
		Prompt:     predict.Input,
		Completion: []string{fmt.Sprintf("%v", resp.Documents[0])},
	}

	for i := range resp.Documents {
		predict.Details.Documents = append(predict.Details.Documents, resp.Documents[i])
	}

	c.AppendToSession(header, body, predict, parameters.CurrentID, modelTrainer)
}

// VisualLogCompletion - Response details
func (c EventManager) VisualLogCompletion(resp *gpt3.CompletionResponse) {
	c.AppendToLayout(c.AppendToChoice(resp, nil, nil, nil))
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

// VisualLogEdit - Log edited response details
func (c EventManager) VisualLogEdit(resp *gpt3.EditsResponse) {
	c.AppendToLayout(c.AppendToChoice(nil, resp, nil, nil))
	node.layout.infoOutput.SetText(fmt.Sprintf("\nCreated: %v\nObject: %v\nCompletion tokens: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
		resp.Created,
		resp.Object,
		resp.Usage.CompletionTokens,
		resp.Usage.PromptTokens,
		resp.Usage.TotalTokens,
		resp.Choices[0].Index))
}

// VisualLogEmbedding - Log embedding response details
func (c EventManager) VisualLogEmbedding(resp *gpt3.EmbeddingsResponse) {
	c.AppendToLayout(c.AppendToChoice(nil, nil, resp, nil))
	node.layout.infoOutput.SetText(fmt.Sprintf("\nObject: %v\nPrompt tokens: %v\nTotal tokens: %v\nIndex: %v\n",
		resp.Object,
		resp.Usage.PromptTokens,
		resp.Usage.TotalTokens,
		resp.Data[0].Index))
}

// VisualLogPredict - Log predicted response details
func (c EventManager) VisualLogPredict(resp *model.PredictResponse) {
	var buffer []string
	for i := range resp.Documents {
		c.AppendToLayout(c.AppendToChoice(nil, nil, nil, &resp.Documents[i]))

		details := fmt.Sprintf("\nAverage probability: %v\nCompletely generated probability: %v\nOverall burstiness: %v",
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

// LogPredictEngine - Log current predict engine
func (c EventManager) LogPredictEngine(client Agent) {
	node.layout.metadataOutput.SetText(
		fmt.Sprintf("\nModel: %v\nAverage Prob: %v\nCompletely Prob: %v\noversall burstiness: %v\n",
			client.engineProperties.Model,
			client.predictProperties.Details.Documents[0].AverageProb,
			client.predictProperties.Details.Documents[0].CompletelyProb,
			client.predictProperties.Details.Documents[0].OverallBurstiness))
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
