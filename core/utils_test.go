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

		_, err := CleanGPTJson[interface{}](llmResponse)

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

func TestSelectRandomRow(t *testing.T) {

	t.Run("Empty dataset", func(t *testing.T) {
		dataset := []RowItem{}
		item, index := SelectRandomRow(dataset)

		if item != nil {
			t.Errorf("Expected nil, got %v", item)
		}

		if index != -1 {
			t.Errorf("Expected -1, got %v", index)
		}
	})

	t.Run("Valid dataset", func(t *testing.T) {
		dataset := []RowItem{
			{"name": "Alice", "age": "30"},
			{"name": "Bob", "age": "25"},
		}

		row, index := SelectRandomRow(dataset)
		if row == nil {
			t.Errorf("Expected a row, got nil")
		}

		if index < 0 || index >= len(dataset) {
			t.Errorf("Expected an index within the dataset range, got %v", index)
		}
	})
}

func TestCountSelectedRows(t *testing.T) {

	t.Run("Returns 0 for empty dataset", func(t *testing.T) {
		dataset := []RowItem{}
		selected := map[int]bool{}

		count := CountSelectedRows(dataset, selected)

		if count != 0 {
			t.Errorf("Expected 0, got %v", count)
		}
	})

	t.Run("Returns correct value for non-empty dataset", func(t *testing.T) {
		dataset := []RowItem{
			{"name": "Alice", "age": "30"},
			{"name": "Bob", "age": "25"},
		}

		selected := map[int]bool{
			0: true,
			1: true,
		}

		count := CountSelectedRows(dataset, selected)
		if count != 2 {
			t.Errorf("Expected 2, got %v", count)
		}
	})
}
