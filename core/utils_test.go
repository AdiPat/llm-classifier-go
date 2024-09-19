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

func TestReadCSVFile(t *testing.T) {

	t.Run("Empty file path", func(t *testing.T) {
		_, err := ReadCSVFile("")
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("Invalid file path", func(t *testing.T) {
		_, err := ReadCSVFile("invalid")
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("Valid CSV file", func(t *testing.T) {
		filePath := "../datasets/student_performance.csv"
		_, err := ReadCSVFile(filePath)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
