package core

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Label = string
type LabelDescription = string

type TaoClassifier struct {
	prompts      map[Label][]LabelDescription // mapping between label -> array of descriptions of label (prompts)
	openaiClient *openai.Client
}

type ClassificationResult struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}

func NewTaoClassifier() *TaoClassifier {
	err := godotenv.Load("../.env")

	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	var openAiKey string = os.Getenv("OPENAI_API_KEY")

	if openAiKey == "" {
		panic("OPENAI_KEY is not set")
	}

	client := openai.NewClient(option.WithAPIKey(openAiKey))

	return &TaoClassifier{
		prompts:      make(map[Label][]LabelDescription),
		openaiClient: client,
	}
}

func (c *TaoClassifier) PromptTrain(prompts map[Label][]LabelDescription) (bool, error) {
	if len(prompts) == 0 {
		return false, fmt.Errorf("prompts cannot be empty")
	}

	for label, descriptionList := range prompts {
		c.prompts[label] = append(c.prompts[label], descriptionList...)
	}

	return true, nil
}

func (c *TaoClassifier) AddPrompt(label Label, description LabelDescription) (bool, error) {
	if label == "" {
		return false, fmt.Errorf("label cannot be empty")
	}

	if description == "" {
		return false, fmt.Errorf("description cannot be empty")
	}

	if _, ok := c.prompts[label]; !ok {
		c.prompts[label] = []LabelDescription{description}
	} else {
		c.prompts[label] = append(c.prompts[label], description)
	}

	return true, nil
}

func (c *TaoClassifier) GetPrompt(label Label) ([]LabelDescription, error) {
	if label == "" {
		return []LabelDescription{}, fmt.Errorf("label cannot be empty")
	}

	descriptionList, ok := c.prompts[label]

	if !ok {
		return []LabelDescription{}, fmt.Errorf("label not found")
	}

	return descriptionList, nil
}

func (c *TaoClassifier) GetPrompts() map[Label][]LabelDescription {
	return c.prompts
}

func (c *TaoClassifier) RemovePrompt(label Label) (bool, error) {
	if label == "" {
		return false, fmt.Errorf("label cannot be empty")
	}

	_, ok := c.prompts[label]

	if !ok {
		return false, fmt.Errorf("label not found")
	}

	delete(c.prompts, label)
	return true, nil
}

func (c *TaoClassifier) ClearPrompts() {
	c.prompts = make(map[Label][]LabelDescription)
}

func (c *TaoClassifier) PredictOne(text string) (ClassificationResult, error) {
	ctx := context.Background()

	// convert c.prompts to a string
	labelDescriptors := "Labels->Description\n"
	for label, descriptionList := range c.prompts {
		for _, description := range descriptionList {
			labelDescriptors += fmt.Sprintf("%s: %s\n", label, description)
		}
	}

	// TODO: Implement OpenAI API call
	systemPrompt := fmt.Sprintf(`You are an AI assistant that performs classification. 
	You will be given a map of labels and their corresponding descriptions. 
	Use this information to classify the given data point.
	Respond in JSON with { label: <label>, "probability": <probability> }. 
	The label should be only from the given labels.
	Context: %s\n`, labelDescriptors)

	if text == "" {
		return ClassificationResult{Label: "", Probability: -1}, fmt.Errorf("text cannot be empty")
	}

	userPrompt := fmt.Sprintf(`Classify the following text: "%s"`, text)

	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT4oMini),
	}

	completions, err := c.openaiClient.Chat.Completions.New(ctx, params)

	if err != nil {
		return ClassificationResult{Label: "", Probability: -1}, err
	}

	result, err := CleanGPTJson[ClassificationResult](completions.Choices[0].Message.Content)

	if err != nil {
		return ClassificationResult{Label: "", Probability: -1}, err
	}

	return result, nil
}

func (c *TaoClassifier) PredictMany(texts []string) ([]ClassificationResult, error) {
	if len(texts) == 0 {
		return []ClassificationResult{}, fmt.Errorf("texts cannot be empty")
	}

	var results []ClassificationResult

	for _, text := range texts {
		// TODO: if multiple values fit into the prompt, then use a single prompt - this is an area of optimization
		result, err := c.PredictOne(text)

		if err != nil {
			return []ClassificationResult{}, err
		}

		results = append(results, result)
	}

	return results, nil
}
