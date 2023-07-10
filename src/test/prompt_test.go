// Test section - Use case
package caos

import (
	"testing"

	"caos/model"
	"caos/service"
)

func TestSendChatCompletion(t *testing.T) {
	t.Run("SendChatCompletion", func(t *testing.T) {
		controller := &service.Controller{}
		prompt := &service.Prompt{}

		agent := controller.AttachProfile()

		engineProperties := &model.EngineProperties{
			UserID:           "test_user",
			Model:            "gpt-3.5-turbo",
			Role:             model.Assistant,
			Temperature:      1.0,
			TopP:             0.4,
			PresencePenalty:  0.5,
			FrequencyPenalty: 0.5,
		}

		promptProperties := &model.PromptProperties{
			Input:         []string{"Generate an uml template"},
			Instruction:   []string{"for an eshop, include customers and providers."},
			Content:       []string{""},
			MaxTokens:     1024,
			Results:       1,
			Probabilities: 1,
		}

		templateProperties := &model.TemplateProperties{}

		agent.EngineProperties = *engineProperties
		agent.PromptProperties = *promptProperties
		agent.TemplateProperties = *templateProperties

		chat, schat := prompt.SendChatCompletionPrompt(agent)
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
		controller := &service.Controller{}
		prompt := &service.Prompt{}

		agent := controller.AttachProfile()

		engineProperties := &model.EngineProperties{
			UserID:           "test_user",
			Model:            "text-davinci-003",
			Role:             model.Assistant,
			Temperature:      1.0,
			TopP:             0.4,
			PresencePenalty:  0.5,
			FrequencyPenalty: 0.5,
		}

		promptProperties := &model.PromptProperties{
			Input:         []string{"Generate an uml template"},
			Instruction:   []string{"for an eshop, include customers and providers."},
			Content:       []string{""},
			MaxTokens:     1024,
			Results:       1,
			Probabilities: 1,
		}

		templateProperties := &model.TemplateProperties{}

		agent.EngineProperties = *engineProperties
		agent.PromptProperties = *promptProperties
		agent.TemplateProperties = *templateProperties

		chat := prompt.SendCompletionPrompt(agent)
		if chat != nil {
			t.Log("Test - PASSED")
		} else {
			t.Errorf("Received:%v", chat)
			t.Log("Test - FAILED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestSendEditPrompt(t *testing.T) {
	t.Run("SendEditPrompt", func(t *testing.T) {
		controller := &service.Controller{}
		prompt := &service.Prompt{}

		agent := controller.AttachProfile()

		engineProperties := &model.EngineProperties{
			UserID:           "test_user",
			Model:            "text-davinci-edit-001",
			Role:             model.Assistant,
			Temperature:      1.0,
			TopP:             0.4,
			PresencePenalty:  0.5,
			FrequencyPenalty: 0.5,
		}

		promptProperties := &model.PromptProperties{
			Input:         []string{"Generate an uml template"},
			Instruction:   []string{"for an eshop, include customers and providers."},
			Content:       []string{"Instructive content"},
			MaxTokens:     512,
			Results:       1,
			Probabilities: 1,
		}

		agent.EngineProperties = *engineProperties
		agent.PromptProperties = *promptProperties

		resp := prompt.SendEditPrompt(agent)
		if resp == nil {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}

		t.Log("Test - FINISHED")
	})
}

func TestSendEmbeddingPrompt(t *testing.T) {
	t.Run("SendEmbeddingPrompt", func(t *testing.T) {
		controller := &service.Controller{}
		prompt := &service.Prompt{}

		agent := controller.AttachProfile()

		engineProperties := &model.EngineProperties{
			UserID:           "test_user",
			Model:            "text-embedding-ada-002",
			Role:             model.Assistant,
			Temperature:      1.0,
			TopP:             0.4,
			PresencePenalty:  0.5,
			FrequencyPenalty: 0.5,
		}

		promptProperties := &model.PromptProperties{
			Input:         []string{"Generate an uml template"},
			Instruction:   []string{"for an eshop, include customers and providers."},
			Content:       []string{"Instructive content"},
			MaxTokens:     1024,
			Results:       1,
			Probabilities: 1,
		}

		agent.EngineProperties = *engineProperties
		agent.PromptProperties = *promptProperties

		resp := prompt.SendEmbeddingPrompt(agent)
		if resp == nil {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestSendPredictablePrompt(t *testing.T) {
	t.Run("SendPredictablePrompt", func(t *testing.T) {
		controller := &service.Controller{}
		prompt := &service.Prompt{}

		agent := controller.AttachProfile()

		predictProperties := model.PredictProperties{
			Input:   []string{"This is a generative response"},
			Details: model.PredictResponse{},
		}

		agent.PredictProperties = predictProperties

		resp := prompt.SendEmbeddingPrompt(agent)
		if resp == nil {
			t.Errorf("Received:%v", resp)
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestGetListModels(t *testing.T) {
	t.Run("GetListModels", func(t *testing.T) {
		controller := &service.Controller{}
		prompt := &service.Prompt{}

		agent := controller.AttachProfile()
		engines := prompt.GetListModels(agent)

		if engines == nil {
			t.Errorf("Received:%v", engines)
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}

		t.Log("Test - FINISHED")
	})
}
