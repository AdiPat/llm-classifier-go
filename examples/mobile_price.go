package main

import (
	"fmt"

	LLMClassifier "github.com/AdiPat/llm-classifier-go/core"
)

func main() {
	params := LLMClassifier.TaoClassifierOptions{
		TrainingDatasetPath: "./datasets/mobile_price_train.csv",
		TargetColumn:        "price_range",
		Verbose:             false,
		PromptSampleSize:    2,
	}

	fmt.Println("Initializing the classifier. ")
	classifier := LLMClassifier.NewTaoClassifier(params)

	fmt.Println("Training the classifier. This may take a while. ")
	classifier.Train()

	testSet, err := LLMClassifier.ReadCSVFile("./datasets/mobile_price_test.csv")

	if err != nil {
		println("Error: Failed to read test dataset", err)
	}

	// TODO: This should not be needed. The classifier should be able to handle []RowItem
	// Convert testSet to []any

	var testSetAny []any
	for _, item := range testSet {
		testSetAny = append(testSetAny, item)
	}

	// use only 10 entries from the test set
	testSetAny = testSetAny[:10]

	fmt.Println("Test set: ", testSetAny)

	predictions, err := classifier.PredictManyObjects(testSetAny)

	if err != nil {
		println("PredictManyObjects failed. Error:", err)
	}

	for index, prediction := range predictions {
		fmt.Println("RowItem: ", testSet[index])
		fmt.Println("Prediction: ", prediction)
	}

}
