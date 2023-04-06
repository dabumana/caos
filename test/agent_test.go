// Test section - Use case
package caos

import (
	"testing"

	"caos/model"
	"caos/service"
)

func TestConnect(t *testing.T) {
	var service service.Agent
	client, exClient := service.Connect()
	if client == nil || exClient == nil {
		t.Error("client not found.")
		t.Log("Test - ERROR")
	} else {
		t.Log("Test - PASSED")
	}
	t.Log("Test - FINISHED")
}

func TestSetEngineParameters(t *testing.T) {
	var service service.Agent
	id := "test_user"
	role := model.Assistant
	model := "text-davinci-003"
	temperature := 1.0
	topp := 0.4
	penalty := 0.5
	frequency := 0.5
	engineProperties := service.SetEngineParameters(
		id,
		model,
		role,
		float32(temperature),
		float32(topp),
		float32(penalty),
		float32(frequency))
	if engineProperties.Model != model ||
		engineProperties.Temperature != float32(temperature) ||
		engineProperties.TopP != float32(topp) ||
		engineProperties.PresencePenalty != float32(penalty) ||
		engineProperties.FrequencyPenalty != float32(frequency) {
		t.Errorf("Received:%v\nExpected:%v\n", engineProperties.Model, model)
		t.Errorf("Received:%v\nExpected:%v\n", engineProperties.Temperature, temperature)
		t.Errorf("Received:%v\nExpected:%v\n", engineProperties.TopP, topp)
		t.Errorf("Received:%v\nExpected:%v\n", engineProperties.PresencePenalty, penalty)
		t.Errorf("Received:%v\nExpected:%v\n", engineProperties.FrequencyPenalty, frequency)
		t.Log("Test - ERROR")
	} else {
		t.Log("Test - PASSED")
	}
	t.Log("Test - FINISHED")
}

func TestSetPromptParameters(t *testing.T) {
	var service service.Agent
	context := []string{"Generate an UML template"}
	prompt := []string{"for an eshop, include customers and providers."}
	tokens := 64
	result := 4
	probabilities := 4
	requestProperties := service.SetPromptParameters(
		context,
		prompt,
		tokens,
		result,
		probabilities)
	if requestProperties.PromptContext == nil ||
		requestProperties.Instruction == nil ||
		requestProperties.MaxTokens != tokens ||
		requestProperties.Results != result ||
		requestProperties.Probabilities != probabilities {
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.PromptContext, context)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Instruction, prompt)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.MaxTokens, tokens)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Results, result)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Probabilities, probabilities)
		t.Log("Test - ERROR")
	} else {
		t.Log("Test - PASSED")
	}
	t.Log("Test - FINISHED")
}

func TestSetPredictionParameters(t *testing.T) {
	var service service.Agent
	context := []string{"Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit..."}
	predictProperties := service.SetPredictionParameters(context)
	if predictProperties.Input == nil {
		t.Errorf("Received:%v\nExpected:%v\n", predictProperties.Input, context)
		t.Error("Test - ERROR")
	} else {
		t.Log("Test - PASSED")
	}
	t.Log("Test - FINISHED")
}
