package core

import (
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

func NewTaoClassifier() *TaoClassifier {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	var openAiKey string = os.Getenv("OPENAI_KEY")

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

func (c *TaoClassifier) Classify(text string) (Label, error) {
	// TODO: Implement OpenAI API call

	if text == "" {
		return "", fmt.Errorf("text cannot be empty")
	}

	return "", nil
}
