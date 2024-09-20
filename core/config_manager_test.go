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

func TestDeleteConfig(t *testing.T) {
	t.Run("Deletes config folder correctly. ", func(t *testing.T) {
		config := GetTaoConfig()

		err := config.DeleteConfigFolder()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestSaveModelToFile(t *testing.T) {
	t.Run("Saves model to file correctly. ", func(t *testing.T) {
		config := GetTaoConfig()

		model := SavedTaoModel{
			ModelId: "test_model",
			Prompts: map[Label][]LabelDescription{
				"0": {"low parental support"},
				"1": {"medium parental support"},
				"2": {"high parental support"},
			},
			Temperature:      0.5,
			PromptSampleSize: 10,
			TargetColumn:     "ParentalSupport",
		}

		err := config.SaveModelToFile(model, SaveModelToFileOptions{Overwrite: true})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
