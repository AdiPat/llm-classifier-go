package core

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AI struct {
	openaiClient *openai.Client
}

type GenerateTextOptions struct {
	System      string
	Temperature float64
	Verbose     bool
}

func NewAI() *AI {
	var openAiKey string = os.Getenv("OPENAI_API_KEY")

	if openAiKey == "" {
		panic("OPENAI_KEY is not set. Please set the OPENAI_API_KEY environment variable.")
	}

	client := openai.NewClient(option.WithAPIKey(openAiKey))

	return &AI{
		openaiClient: client,
	}
}

func (ai *AI) GenerateText(prompt string, opts ...GenerateTextOptions) (string, error) {
	options := GenerateTextOptions{
		Temperature: 0.5, // Default temperature
		System:      "You are a helpful AI-assistant that generates text based on the given prompt.",
		Verbose:     false,
	}

	// Override with provided options if any
	if len(opts) > 0 {
		options = opts[0]
	}

	if options.Verbose {
		fmt.Println("Generating text with prompt:", prompt)
	}

	ctx := context.Background()

	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(options.System),
			openai.UserMessage(prompt),
		}),
		Seed:        openai.Int(1),
		Model:       openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.Float(options.Temperature),
	}

	completions, err := ai.openaiClient.Chat.Completions.New(ctx, params)

	if err != nil {
		log.Fatal("GenerateText: Failed to generate completions. ", err)
		return "", err
	}

	result := completions.Choices[0].Message.Content

	if options.Verbose {
		fmt.Println("LLM Response: ", result)
	}

	return result, nil
}

func (ai *AI) GenerateObject(prompt string, schema string, opts ...GenerateTextOptions) (interface{}, error) {
	options := GenerateTextOptions{
		Temperature: 0.5, // Default temperature
		System:      "You are a helpful AI-assistant that generates text based on the given prompt.",
		Verbose:     false,
	}

	promptWithSchema := fmt.Sprintf("%s\n Return the response in JSON as per the schema. \nSchema: %s", prompt, schema)

	// Override with provided options if any
	if len(opts) > 0 {
		options = opts[0]
	}

	if options.Verbose {
		fmt.Println("Generating text with prompt:", promptWithSchema)
	}

	ctx := context.Background()

	params := openai.ChatCompletionNewParams{

		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(options.System),
			openai.UserMessage(promptWithSchema),
		}),
		Seed:        openai.Int(1),
		Model:       openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.Float(options.Temperature),
	}

	completions, err := ai.openaiClient.Chat.Completions.New(ctx, params)

	contentStr := completions.Choices[0].Message.Content

	if options.Verbose {
		fmt.Println("LLM Response: ", contentStr)
	}

	if err != nil {
		log.Fatal("GenerateObject: Failed to generate completions. ", err)
		return nil, err
	}

	result, err := CleanGPTJson[interface{}](contentStr)

	if err != nil {
		log.Fatal("GenerateObject: Failed to generate completions. ", err)
		return nil, err
	}

	return result, nil
}
