package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type TaoConfig struct {
	configFolder string
	modelsFolder string
}

var (
	instance *TaoConfig
	once     sync.Once
)

func newTaoConfig() *TaoConfig {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return nil
	}

	configFolder := filepath.Join(homeDir, ".tao")
	modelsFolder := filepath.Join(configFolder, "models")

	return &TaoConfig{
		configFolder: configFolder,
		modelsFolder: modelsFolder,
	}
}

func (tc *TaoConfig) Init() error {
	err := CreateFolderIfNotExists(tc.configFolder)

	if err != nil {
		fmt.Println("Error creating config folder:", err)
		return err
	}

	err = CreateFolderIfNotExists(tc.modelsFolder)

	if err != nil {
		fmt.Println("Error creating models folder:", err)
		return err
	}

	return nil
}

func (tc *TaoConfig) DeleteConfigFolder() error {
	err := os.RemoveAll(tc.modelsFolder)

	if err != nil {
		fmt.Println("Error deleting models folder:", err)
		return err
	}

	err = os.RemoveAll(tc.configFolder)

	if err != nil {
		fmt.Println("Error deleting config folder:", err)
		return err
	}

	return nil
}

// GetTaoConfig returns the singleton instance of TaoConfig
func GetTaoConfig() *TaoConfig {
	once.Do(func() {
		instance = newTaoConfig()
		instance.Init()
	})
	return instance
}

type SaveModelToFileOptions struct {
	Overwrite bool
}

func (tc *TaoConfig) SaveModelToFile(model SavedTaoModel, opts ...SaveModelToFileOptions) (string, error) {
	modelId := model.ModelId

	if modelId == "" {
		// TODO: not sure if this is a good idea to autogenerate the ID
		// TODO: use UUID instead for safety in distributed systems
		fmt.Println("SaveModelToFile: ModelId is empty. Autogenerating ID. ")
		modelId = fmt.Sprint(time.Now().Unix())
	}

	options := SaveModelToFileOptions{
		Overwrite: true,
	}

	if len(opts) > 0 {
		options = opts[0]
	}

	// Save the model to the file
	modelFilePath := filepath.Join(tc.modelsFolder, modelId+".json")

	taoModelBytes, err := json.Marshal(model)

	if err != nil {
		fmt.Println("SaveModelToFile: Error marshalling model:", err)
		return "", err
	}

	// check if file exists
	if _, err := os.Stat(modelFilePath); err == nil {
		// file exists
		if options.Overwrite {
			fmt.Println("SaveModelToFile: Model file already exists. Rewriting.")
		} else {
			fmt.Println("SaveModelToFile: Model file already exists. Skipping.")
			return "", nil
		}
	}

	err = os.WriteFile(modelFilePath, taoModelBytes, 0644)

	if err != nil {
		fmt.Println("SaveModelToFile: Error writing model to file:", err)
		return "", err
	}

	return modelId, nil
}
