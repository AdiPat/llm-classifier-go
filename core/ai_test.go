package core

import (
	"testing"
)

func TestGenerateText(t *testing.T) {
	t.Run("Generates some text based on arbitrary prompt. ", func(t *testing.T) {
		ai := NewAI()

		result, err := ai.GenerateText("Hello, what's your name?", GenerateTextOptions{Verbose: true})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(result) == 0 {
			t.Errorf("Expected non-empty text, got empty")
		}
	})
}
