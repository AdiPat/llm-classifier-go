package examples

import (
	"fmt"

	LLMClassifier "github.com/AdiPat/llm-classifier-go/core"
)

func MobilePriceExample() {
	params := LLMClassifier.TaoClassifierOptions{
		ModelId:             "mobile_price_classifier",
		TrainingDatasetPath: "./datasets/mobile_price_train.csv",
		TargetColumn:        "price_range",
		Verbose:             false,
		PromptSampleSize:    2,
	}

	fmt.Println("Initializing the classifier. ")
	classifier := LLMClassifier.NewTaoClassifier(params)

	fmt.Println("Loading model if exists.")

	_, err := classifier.LoadModel("mobile_price_classifier")

	if err != nil {
		println("Error: Failed to load model", err)
		fmt.Println("Training the classifier. This may take a while. ")
		classifier.Train()
	}

	fmt.Println("Saving the model. ")

	_, err = classifier.SaveModel()

	if err != nil {
		println("Error: Failed to save model", err)
	}

	testSet, err := LLMClassifier.ReadCSVFile("./datasets/mobile_price_test.csv")

	if err != nil {
		println("Error: Failed to read test dataset", err)
	}

	testSet = testSet[:10]

	fmt.Println("Test set: ", testSet)

	predictions, err := classifier.PredictManyRowItems(testSet)

	if err != nil {
		println("PredictManyObjects failed. Error:", err)
	}

	for index, prediction := range predictions {
		fmt.Println("RowItem: ", testSet[index])
		fmt.Println("Prediction: ", prediction)
	}

}
