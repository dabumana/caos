// Test section - Use case
package caos

import (
	"testing"

	"caos/model"
	"caos/service"
)

func TestExecuteChainJob(t *testing.T) {
	t.Run("ExecuteChainJob", func(t *testing.T) {
		agent := &service.Agent{}
		chain := &service.Chain{}
		prompt := model.PromptProperties{
			Input: []string{"Local time"},
		}

		agent.Initialize()
		chain.ExecuteChainJob(*agent, &prompt)

		if chain.Transform.Source[0] == "" {
			t.Errorf("Received:%v", chain.Transform.Source)
			t.Log("Test - ERROR")
		} else {
			t.Log("Test - PASSED")
		}

		t.Log("Test - FINISHED")
	})
}
