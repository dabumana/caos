// Test section - Use case
package caos

import (
	"testing"

	"caos/model"
	"caos/service"
)

var controller = &service.Controller{}
var prompter = &service.Prompt{}

var localAgent = controller.AttachProfile()

var engineProperties = &model.EngineProperties{
	UserID:           "test_user",
	Role:             model.Assistant,
	Temperature:      1.0,
	TopP:             0.4,
	PresencePenalty:  0.5,
	FrequencyPenalty: 0.5,
}

var promptProperties = &model.PromptProperties{
	Input:         []string{"Generate an uml template"},
	Instruction:   []string{"for an eshop, include customers and providers."},
	Content:       []string{"UML generated"},
	MaxTokens:     1024,
	Results:       1,
	Probabilities: 1,
}

var predictProperties = &model.PredictProperties{
	Input:   []string{"This is a generative response"},
	Details: model.PredictResponse{},
}

var templateProperties = &model.TemplateProperties{}

func initializeAgent() {
	localAgent.EngineProperties = *engineProperties
	localAgent.PromptProperties = *promptProperties
	localAgent.PredictProperties = *predictProperties
	localAgent.TemplateProperties = *templateProperties
}

func TestSendChatCompletion(t *testing.T) {
	t.Run("SendChatCompletion", func(t *testing.T) {
		engineProperties.Model = "gpt-3.5-turbo"
		initializeAgent()

		chat, schat := prompter.SendChatCompletionPrompt(localAgent)
		if chat != nil && schat == nil ||
			schat != nil && chat == nil {
			t.Log("Test - PASSED")
		} else {
			t.Errorf("Received:%v/%v", chat, schat)
			t.Log("Test - FAILED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestSendCompletion(t *testing.T) {
	t.Run("SendCompletion", func(t *testing.T) {
		engineProperties.Model = "text-davinci-003"
		initializeAgent()

		resp := prompter.SendCompletionPrompt(localAgent)
		if resp != nil {
			t.Log("Test - PASSED")
		} else {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestSendEditPrompt(t *testing.T) {
	t.Run("SendEditPrompt", func(t *testing.T) {
		engineProperties.Model = "text-davinci-edit-001"
		initializeAgent()

		resp := prompter.SendEditPrompt(localAgent)
		if resp != nil {
			t.Log("Test - PASSED")
		} else {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		}

		t.Log("Test - FINISHED")
	})
}

func TestSendEmbeddingPrompt(t *testing.T) {
	t.Run("SendEmbeddingPrompt", func(t *testing.T) {
		engineProperties.Model = "text-embedding-ada-002"
		initializeAgent()

		resp := prompter.SendEmbeddingPrompt(localAgent)
		if resp != nil {
			t.Log("Test - PASSED")
		} else {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestSendPredictablePrompt(t *testing.T) {
	t.Run("SendPredictablePrompt", func(t *testing.T) {
		initializeAgent()

		resp := prompter.SendPredictablePrompt(localAgent)
		if resp != nil {
			t.Log("Test - PASSED")
		} else {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestGetListModels(t *testing.T) {
	t.Run("GetListModels", func(t *testing.T) {
		initializeAgent()

		resp := prompter.GetListModels(localAgent)
		if resp != nil {
			t.Log("Test - PASSED")
		} else {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		}
		t.Log("Test - FINISHED")
	})
}
