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

func checkTest(t *testing.T, resp any) {
	if resp != nil {
		t.Log("Test - PASSED")
	} else {
		t.Errorf("Received:%v", resp)
		t.Log("Test - FAILED")
	}
}

func checkBenchmark(t *testing.B, resp any) {
	if resp != nil {
		t.Log("Test - PASSED")
	} else {
		t.Errorf("Received:%v", resp)
		t.Log("Test - FAILED")
	}
}

func TestSendChatCompletion(t *testing.T) {
	t.Run("SendChatCompletion", func(t *testing.T) {
		engineProperties.Model = "gpt-3.5-turbo"
		initializeAgent()

		resp, _ := prompter.SendChatCompletionPrompt(localAgent)
		checkTest(t, resp)
		t.Log("Test - FINISHED")
	})
}

func BenchmarkSendChatCompletion(t *testing.B) {
	t.Run("SendChatCompletion", func(t *testing.B) {
		engineProperties.Model = "gpt-3.5-turbo"
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp, _ := prompter.SendChatCompletionPrompt(localAgent)
			checkBenchmark(t, resp)
		}
		t.Log("Test - FINISHED")
	})
}

func TestSendCompletion(t *testing.T) {
	t.Run("SendCompletion", func(t *testing.T) {
		engineProperties.Model = "text-davinci-003"
		initializeAgent()

		resp := prompter.SendCompletionPrompt(localAgent)
		checkTest(t, resp)
	})
	t.Log("Test - FINISHED")
}

func BenchmarkSendCompletion(t *testing.B) {
	t.Run("SendCompletion", func(t *testing.B) {
		engineProperties.Model = "text-davinci-003"
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp := prompter.SendCompletionPrompt(localAgent)
			checkBenchmark(t, resp)
		}
	})
	t.Log("Test - FINISHED")
}

func TestSendEditPrompt(t *testing.T) {
	t.Run("SendEditPrompt", func(t *testing.T) {
		engineProperties.Model = "text-davinci-edit-001"
		initializeAgent()

		resp := prompter.SendEditPrompt(localAgent)
		checkTest(t, resp)
	})
	t.Log("Test - FINISHED")
}

func BenchmarkSendEditPrompt(t *testing.B) {
	t.Run("SendCompletion", func(t *testing.B) {
		engineProperties.Model = "text-davinci-edit-001"
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp := prompter.SendEditPrompt(localAgent)
			checkBenchmark(t, resp)
		}
	})
	t.Log("Test - FINISHED")
}

func TestSendEmbeddingPrompt(t *testing.T) {
	t.Run("SendEmbeddingPrompt", func(t *testing.T) {
		engineProperties.Model = "text-embedding-ada-002"
		initializeAgent()

		resp := prompter.SendEmbeddingPrompt(localAgent)
		checkTest(t, resp)
	})
	t.Log("Test - FINISHED")
}

func BenchmarkSendEmbeddingPrompt(t *testing.B) {
	t.Run("SendEmbeddingPrompt", func(t *testing.B) {
		engineProperties.Model = "text-embedding-ada-002"
		initializeAgent()

		t.ReportAllocs()
		for i := 0; i < t.N; i++ {
			resp := prompter.SendEmbeddingPrompt(localAgent)
			checkBenchmark(t, resp)
		}
	})
	t.Log("Test - FINISHED")
}

func TestSendPredictablePrompt(t *testing.T) {
	t.Run("SendPredictablePrompt", func(t *testing.T) {
		initializeAgent()

		resp := prompter.SendPredictablePrompt(localAgent)
		checkTest(t, resp)
	})
	t.Log("Test - FINISHED")
}

func BenchmarkSendPredictablePrompt(t *testing.B) {
	t.Run("SendPredictablePrompt", func(t *testing.B) {
		initializeAgent()

		for i := 0; i < t.N; i++ {
			resp := prompter.SendPredictablePrompt(localAgent)
			checkBenchmark(t, resp)
		}
	})
	t.Log("Test - FINISHED")
}

func TestGetListModels(t *testing.T) {
	t.Run("GetListModels", func(t *testing.T) {
		initializeAgent()

		resp := prompter.GetListModels(localAgent)
		checkTest(t, resp)
	})
	t.Log("Test - FINISHED")
}

func BenchmarkGetListModels(t *testing.B) {
	t.Run("GetListModels", func(b *testing.B) {
		initializeAgent()

		for i := 0; i < b.N; i++ {
			resp := prompter.GetListModels(localAgent)
			checkBenchmark(t, resp)
		}
	})
	t.Log("Test - FINISHED")
}
