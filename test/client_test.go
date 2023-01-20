// Test section - Use case
package caos

import (
	"caos/service"
	"testing"
)

func TestConnect(t *testing.T) {
	var service service.Client
	client := service.Connect()
	if client == nil {
		t.Error("client not found.")
		t.Log("Test - ERROR")
	} else {
		t.Log("Test - PASSED")
	}
	t.Log("Test - FINISHED")
}

func TestSetEngineParameters(t *testing.T) {
	var service service.Client
	model := "text-davinci-003"
	temperature := 1.0
	topp := 0.4
	penalty := 0.5
	frequency := 0.5
	engineProperties := service.SetEngineParameters(
		model,
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

func TestSetRequestParameters(t *testing.T) {
	var service service.Client
	context := []string{"Generate an UML template"}
	prompt := []string{"for an eshop, include customers and providers."}
	tokens := 64
	result := 4
	probabilities := 4
	requestProperties := service.SetRequestParameters(
		context,
		prompt,
		tokens,
		result,
		probabilities)
	if requestProperties.PromptContext == nil ||
		requestProperties.Prompt == nil ||
		requestProperties.MaxTokens != tokens ||
		requestProperties.Results != result ||
		requestProperties.Probabilities != probabilities {
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.PromptContext, context)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Prompt, prompt)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.MaxTokens, tokens)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Results, result)
		t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Probabilities, probabilities)
		t.Log("Test - ERROR")
	} else {
		t.Log("Test - PASSED")
	}
	t.Log("Test - FINISHED")
}
