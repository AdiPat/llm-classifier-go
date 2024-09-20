package examples

import (
	"fmt"

	LLMClassifier "github.com/AdiPat/llm-classifier-go/core"
)

func TwitterSentimentAnalysisExample() {

	params := LLMClassifier.TaoClassifierOptions{
		ModelId:          "twitter_sentiment_analysis",
		TargetColumn:     "Sentiment",
		PromptSampleSize: 2,
	}

	classifier := LLMClassifier.NewTaoClassifier(params)

	classifier.PromptTrain(map[LLMClassifier.Label][]LLMClassifier.LabelDescription{
		"positive": {"positive sentiment"},
		"neutral":  {"neutral sentiment"},
		"negative": {"negative sentiment"},
	})

	testSet, err := LLMClassifier.ReadCSVFile("./datasets/twitter_validation.csv")

	if err != nil {
		println("Error: Failed to read test dataset", err)
		panic(err)
	}

	testSet = testSet[:10]

	testSetWithoutSentiment := make([]LLMClassifier.RowItem, len(testSet))

	for _, rowItem := range testSet {
		rowItemWithoutSentiment := make(LLMClassifier.RowItem)

		for key, value := range rowItem {
			if key != "Sentiment" {
				rowItemWithoutSentiment[key] = value
			}
		}

		testSetWithoutSentiment = append(testSetWithoutSentiment, rowItemWithoutSentiment)
	}

	predictions, err := classifier.PredictManyRowItems(testSetWithoutSentiment)

	fmt.Println("Test set: ", len(testSetWithoutSentiment))
	fmt.Println("Predictions: ", len(predictions))

	if err != nil {
		println("PredictManyObjects failed. Error:", err)
		panic(err)
	}

	for index, prediction := range predictions {
		println("RowItem: ", testSetWithoutSentiment[index])
		fmt.Printf("Prediction: %+v\n", prediction)
	}

}
