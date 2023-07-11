// Test section - Use case
package caos

import (
	"testing"

	"caos/model"
	"caos/service"
)

var agent = &service.Agent{}

var context = []string{"Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit..."}
var prompt = []string{"Extend the quote"}

const id = "test_user"
const role = model.Assistant
const temperature = 1.0
const topp = 0.4
const penalty = 0.5
const frequency = 0.5
const result = 4
const probabilities = 4

func TestInitialize(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		agent.Initialize()
		preferences := agent.GetStatus()

		if preferences.TemplateIDs == 0 {
			t.Errorf("Received:%v", preferences.TemplateIDs)
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}

		t.Log("Test - FINISHED")
	})
}

func TestSetEngineParameters(t *testing.T) {
	t.Run("SetEngineParameters", func(t *testing.T) {
		model := "text-davinci-003"
		engineProperties := agent.SetEngineParameters(
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
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}
	})
	t.Log("Test - FINISHED")
}

func TestSetPromptParameters(t *testing.T) {
	t.Run("SetPromptParameters", func(t *testing.T) {
		requestProperties := agent.SetPromptParameters(
			context,
			prompt,
			result,
			probabilities)

		if requestProperties.Input == nil ||
			requestProperties.Instruction == nil ||
			requestProperties.Results != result ||
			requestProperties.Probabilities != probabilities {
			t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Input, context)
			t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Instruction, prompt)
			t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Results, result)
			t.Errorf("Received:%v\nExpected:%v\n", requestProperties.Probabilities, probabilities)
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}
		t.Log("Test - FINISHED")
	})
}

func TestSetPredictionParameters(t *testing.T) {
	t.Run("SetPredictionParameters", func(t *testing.T) {
		predictProperties := agent.SetPredictionParameters(context)
		if predictProperties.Input == nil {
			t.Errorf("Received:%v\nExpected:%v\n", predictProperties.Input, context)
			t.Error("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}
	})
	t.Log("Test - FINISHED")
}
