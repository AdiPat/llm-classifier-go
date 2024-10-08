package core

import (
	"testing"
)

func TestGenerateText(t *testing.T) {
	t.Run("Generates some text based on arbitrary prompt. ", func(t *testing.T) {
		ai := NewAI()

		result, err := ai.GenerateText("Hello, what's your name?", GenerateTextOptions{Verbose: false})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(result) == 0 {
			t.Errorf("Expected non-empty text, got empty")
		}
	})
}

func TestGenerateObject(t *testing.T) {
	t.Run("Generates an object based on the provided schema. ", func(t *testing.T) {
		ai := NewAI()

		schema := `{
			"name": "string",
			"age": "int",
			"description": "string",
			"attributes": "map[string]string",
			"tags": "[]string"
		}`

		_, err := ai.GenerateObject("Generate an object.", schema, GenerateTextOptions{Verbose: false})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

	})
}
