package core

import (
	"testing"
)

func TestPredictOne(t *testing.T) {

	t.Run("Correctly classifies single example", func(t *testing.T) {
		classifier := NewTaoClassifier()

		prompts := map[Label][]LabelDescription{
			"cat": {"Cats are generally more independent and aloof than dogs, who are often more social and affectionate. Cats are also more territorial and may be more aggressive when defending their territory.  Cats are self-grooming animals, using their tongues to keep their coats clean and healthy. Cats use body language and vocalizations, such as meowing and purring, to communicate."},
			"dog": {"Dogs are more pack-oriented and tend to be more loyal to their human family.  Dogs, on the other hand, often require regular grooming from their owners, including brushing and bathing. Dogs use body language and barking to convey their messages. Dogs are also more responsive to human commands and can be trained to perform a wide range of tasks."},
		}

		classifier.PromptTrain(prompts)

		result, err := classifier.PredictOne("Meow")

		if result.Label != "cat" {
			t.Errorf("Expected cat, got %v", result.Label)
		}

		if result.Probability < 0.5 {
			t.Errorf("Expected probability > 0.5, got %v", result.Probability)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		resultDog, err := classifier.PredictOne("Woof")

		if resultDog.Label != "dog" {
			t.Errorf("Expected dog, got %v", resultDog.Label)
		}

		if resultDog.Probability < 0.5 {
			t.Errorf("Expected probability > 0.5, got %v", resultDog.Probability)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestGenerateClassificationProfile(t *testing.T) {

	t.Run("Generates a classifier profile with correct schema", func(t *testing.T) {

		classifier := NewTaoClassifier()

		rowItem := RowItem{
			"student_id":                 "1",
			"name":                       "John",
			"gender":                     "Male",
			"attendance_rate":            "85",
			"study_hours_per_week":       "15",
			"previous_grade":             "78",
			"extracurricular_activities": "1",
			"parental_support":           "High",
			"final_grade":                "80",
		}

		classifierProfile, _ := classifier.GenerateClassifierProfile("parental_support", rowItem, ClassifierProfile{})

		if classifierProfile.Label != "parental_support" {
			t.Errorf("Expected parental_support, got %v", classifierProfile.Label)
		}

		if len(classifierProfile.Description) == 0 {
			t.Errorf("Expected non-empty description, got empty")
		}
	})
}
