// Test section - Use case
package caos

import (
	"testing"

	"caos/service"
)

func TestAttachProfile(t *testing.T) {
	t.Run("AttachProfile", func(t *testing.T) {
		controller := &service.Controller{}
		agent := controller.AttachProfile()
		if agent.GetStatus().User == "" {
			t.Errorf("Received:%v", agent.GetStatus().User)
			t.Log("Test - FAILED")
		} else {
			t.Log("Test - PASSED")
		}
		t.Log("Test - FINISHED")
	})
}
