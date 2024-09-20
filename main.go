package main

import (
	"log"

	"github.com/AdiPat/llm-classifier-go/examples"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// UNCOMMENT EXAMPLES AS NEEDED

	// examples.MobilePriceExample()
	examples.TwitterSentimentAnalysisExample()
}
