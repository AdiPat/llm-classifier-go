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

	t.Run("Strips off plain-text before and after JSON inside markdown", func(t *testing.T) {
		llmResponse := `LLM Response:  Here is a JSON object based on the provided schema:` + "```" + "json" +
			`{
						"name": "Alice Johnson",
						"age": 30,
						"description": "A passionate software developer with a love for open-source projects.",
						"attributes": {
							"key": "innovative thinker"
						},
						"tags": ["developer", "open-source", "technology", "innovation"]
						}` + "```"

		_, err := CleanGPTJson[TestType](llmResponse)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
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
		records, err := ReadCSVFile(filePath)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(records) != 10 {
			t.Errorf("Expected 10 rows found in the CSV file, got %v", len(records))
		}
	})
}
