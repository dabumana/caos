// Test section - Use case
package caos

import (
	"fmt"
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

func checkResponse(t *testing.T, resp any) {
	if resp != nil {
		t.Log("Test - PASSED")
	} else {
		t.Errorf("Received:%v", resp)
		t.Log("Test - FAILED")
	}
	t.Log("Test - FINISHED")
}

func TestSendChatCompletion(t *testing.T) {
	t.Run("SendChatCompletion", func(t *testing.T) {
		engineProperties.Model = "gpt-3.5-turbo"
		initializeAgent()

		resp, _ := prompter.SendChatCompletionPrompt(localAgent)
		checkResponse(t, resp)
	})
}

func BenchmarkSendChatCompletion(t *testing.B) {
	t.Run("SendChatCompletion", func(t *testing.B) {
		engineProperties.Model = "gpt-3.5-turbo"
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp, _ := prompter.SendChatCompletionPrompt(localAgent)
			fmt.Printf("resp: %v\n", resp)
		}
	})
}

func TestSendCompletion(t *testing.T) {
	t.Run("SendCompletion", func(t *testing.T) {
		engineProperties.Model = "text-davinci-003"
		initializeAgent()

		resp := prompter.SendCompletionPrompt(localAgent)
		checkResponse(t, resp)
	})
}

func BenchmarkSendCompletion(t *testing.B) {
	t.Run("SendCompletion", func(t *testing.B) {
		engineProperties.Model = "text-davinci-003"
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp := prompter.SendCompletionPrompt(localAgent)
			fmt.Printf("resp: %v\n", resp)
		}
	})
}

func TestSendEditPrompt(t *testing.T) {
	t.Run("SendEditPrompt", func(t *testing.T) {
		engineProperties.Model = "text-davinci-edit-001"
		initializeAgent()

		resp := prompter.SendEditPrompt(localAgent)
		checkResponse(t, resp)
	})
}

func BenchmarkSendEditPrompt(t *testing.B) {
	t.Run("SendCompletion", func(t *testing.B) {
		engineProperties.Model = "text-davinci-edit-001"
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp := prompter.SendEditPrompt(localAgent)
			fmt.Printf("resp: %v\n", resp)
		}
	})
}

func TestSendEmbeddingPrompt(t *testing.T) {
	t.Run("SendEmbeddingPrompt", func(t *testing.T) {
		engineProperties.Model = "text-embedding-ada-002"
		initializeAgent()

		resp := prompter.SendEmbeddingPrompt(localAgent)
		checkResponse(t, resp)
	})
}

func BenchmarkSendEmbeddingPrompt(t *testing.B) {
	t.Run("SendEmbeddingPrompt", func(t *testing.B) {
		engineProperties.Model = "text-embedding-ada-002"
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp := prompter.SendEmbeddingPrompt(localAgent)
			fmt.Printf("resp: %v\n", resp)
		}
	})
}

func TestSendPredictablePrompt(t *testing.T) {
	t.Run("SendPredictablePrompt", func(t *testing.T) {
		initializeAgent()

		resp := prompter.SendPredictablePrompt(localAgent)
		checkResponse(t, resp)
	})
}

func BenchmarkSendPredictablePrompt(t *testing.B) {
	t.Run("SendPredictablePrompt", func(t *testing.B) {
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp := prompter.SendPredictablePrompt(localAgent)
			fmt.Printf("resp: %v\n", resp)
		}
	})
}

func TestGetListModels(t *testing.T) {
	t.Run("GetListModels", func(t *testing.T) {
		initializeAgent()

		resp := prompter.GetListModels(localAgent)
		checkResponse(t, resp)
	})
}

func BenchmarkGetListModels(t *testing.B) {
	t.Run("GetListModels", func(b *testing.B) {
		initializeAgent()

		for i := 0; i < b.N; i++ {
			prompter.GetListModels(localAgent)
		}
	})
}
