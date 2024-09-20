package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
