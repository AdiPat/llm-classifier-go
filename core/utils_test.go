package core

import (
	"reflect"
	"testing"
)

func TestCleanGPTJSON(t *testing.T) {
	t.Run("Empty JSON string", func(t *testing.T) {
		_, err := CleanGPTJson[interface{}]("")
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("Invalid JSON string", func(t *testing.T) {
		_, err := CleanGPTJson[interface{}]("invalid")
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("Valid JSON string", func(t *testing.T) {
		jsonStr := `{"key": "value"}`
		result, err := CleanGPTJson[interface{}](jsonStr)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := map[string]interface{}{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Valid JSON string inside markdown", func(t *testing.T) {
		jsonStr := "```json\n{\"key\": \"value\"}\n```"
		result, err := CleanGPTJson[interface{}](jsonStr)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := map[string]interface{}{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
