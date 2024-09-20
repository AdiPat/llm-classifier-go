package core

import "testing"

func TestInitConfig(t *testing.T) {

	t.Run("Initializes config correctly. ", func(t *testing.T) {
		config := GetTaoConfig()

		err := config.DeleteConfigFolder()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		err = config.Init()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if config == nil {
			t.Errorf("Expected a non-nil object, got nil")
		}
	})

}
