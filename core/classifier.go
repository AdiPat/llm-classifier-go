package core

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
)

type Class = string
type Label = string
type LabelDescription = string

type TaoClassifier struct {
	prompts          map[Label][]LabelDescription // mapping between label -> array of descriptions of label (prompts)
	ai               *AI
	dataset          []RowItem
	targetColumn     string
	temperature      float64
	promptSampleSize int
	verbose          bool
}

type ClassificationResult struct {
	Label          Label   `json:"label"`
	PredictedClass Class   `json:"predicted_class"`
	Probability    float64 `json:"probability"`
}

type ClassifierProfile struct {
	Label       string   `json:"label"`
	Description []string `json:"description"`
}

type TaoClassifierOptions struct {
	TrainingDatasetPath string
	TargetColumn        string
	Temperature         float64
	PromptSampleSize    int
	Verbose             bool
}

func NewTaoClassifier(opts ...TaoClassifierOptions) *TaoClassifier {
	options := TaoClassifierOptions{
		TrainingDatasetPath: "",
		TargetColumn:        "",
		Temperature:         0.5,
		PromptSampleSize:    10,
	}

	if len(opts) > 0 {
		options = opts[0]
	}

	if options.PromptSampleSize <= 0 {
		// set to default
		options.PromptSampleSize = 10
	}

	err := godotenv.Load("../.env")

	if err != nil {
		if options.Verbose {
			fmt.Println("Error loading .env file", err)
		}
	}

	ai := NewAI()

	dataset := []RowItem{}

	if options.TrainingDatasetPath != "" && options.TargetColumn == "" {
		log.Fatal("NewTaoClassifier: TargetColumn cannot be empty when TrainingDatasetPath is specified. ")
		panic("TargetColumn cannot be empty when TrainingDatasetPath is specified. ")
	}

	prompts := make(map[Label][]LabelDescription)

	if options.TrainingDatasetPath != "" && options.TargetColumn != "" {
		dataset, err = ReadCSVFile(options.TrainingDatasetPath)

		// extract classes from dataset based on target column
		classes := ExtractClasses(dataset, options.TargetColumn)

		if len(classes) == 0 {

			if options.Verbose {
				fmt.Printf("NewTaoClassifier: no classes found on targetColumn=%s for dataset\n", options.TargetColumn)
			}
		}

		for _, class := range classes {
			if _, ok := prompts[class]; !ok {
				prompts[class] = []LabelDescription{}
			}
		}

		if err != nil {
			log.Fatal("NewTaoClassifier: Failed to read training dataset", err)
			panic("NewTaoClassifier: Failed to read training dataset")
		}
	}

	return &TaoClassifier{
		prompts:          prompts,
		ai:               ai,
		dataset:          dataset,
		temperature:      options.Temperature,
		promptSampleSize: options.PromptSampleSize,
		targetColumn:     options.TargetColumn,
		verbose:          options.Verbose,
	}
}

func (c *TaoClassifier) InitializePromptsFromDataset() {
	classes := ExtractClasses(c.dataset, c.targetColumn)

	if len(classes) == 0 {

		if c.verbose {
			fmt.Printf("NewTaoClassifier: no classes found on targetColumn=%s for dataset\n", c.targetColumn)
		}
	}

	for _, class := range classes {
		if _, ok := c.prompts[class]; !ok {
			c.prompts[class] = []LabelDescription{}
		}
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

func (c *TaoClassifier) GetAvailableLabels() ([]Label, error) {
	labels := []Label{}

	for label := range c.prompts {
		labels = append(labels, label)
	}

	return labels, nil
}

func (c *TaoClassifier) GenerateClassifierProfile(label Label, rowItem RowItem, currentClassifierProfile ClassifierProfile) (ClassifierProfile, error) {
	if len(rowItem) == 0 {
		return ClassifierProfile{}, fmt.Errorf("rowItem cannot be empty")
	}

	combinedRowItems := ""

	for key, value := range rowItem {
		combinedRowItems += fmt.Sprintf("%s: %s\n", key, value)
	}

	availableLabels, err := c.GetAvailableLabels()

	if err != nil {
		log.Fatal("GenerateClassifierProfile: Failed to get available labels", err)
		return ClassifierProfile{}, err
	}

	labelsStr := strings.Join(availableLabels, ", ")

	// TODO: find a way to sync the type and schema in the prompt
	systemPrompt := `You are an AI assistant that performs classification.
					You are tasked with generating a 'classification profile' given a set of row items for the given label.
					In the description, include relationships between variables in the row.
					Don't include the row item values in the attributes, include anything additional discovered in the data.
					Respond in JSON with { label: string <label>, "description": string[] <description array> } }.
					Based on the label, identify features within the row items that are relevant to the label.
					Target Column for Classification: ` + c.targetColumn + "\nAvailable Labels: " + labelsStr

	userPrompt := fmt.Sprintf(`Generate a classification profile for the label %s given the following row items: %s`, label, combinedRowItems)

	text, err := c.ai.GenerateText(userPrompt, GenerateTextOptions{Verbose: false, System: systemPrompt})

	if c.verbose {
		fmt.Println("System Prompt: ", systemPrompt)
		fmt.Println("Prompt: ", userPrompt)
		fmt.Println("Generated Text: ", text)
	}

	if err != nil {
		log.Fatal("GenerateClassifierProfile: Failed to generate completions", err)
		return ClassifierProfile{}, err
	}

	result, err := CleanGPTJson[ClassifierProfile](text)

	if err != nil {
		log.Fatal("GenerateClassifierProfile: Failed to generate completions", err)
		return ClassifierProfile{}, err
	}

	return result, nil
}

func (c *TaoClassifier) Train() error {
	maxDescriptions := c.promptSampleSize
	selectedRows := make(map[int]bool)

	for index := range len(c.dataset) {
		selectedRows[index] = false
	}

	if c.verbose {
		fmt.Printf("selectedRows: %+v\n", selectedRows)
		fmt.Println("Dataset Size: ", len(c.dataset))
		fmt.Println("Prompts Before: ", c.prompts)
	}

	c.InitializePromptsFromDataset()

	for {
		row, index := SelectRandomRow(c.dataset)
		selectedRowsCount := CountSelectedRows(c.dataset, selectedRows)

		if selectedRowsCount == len(c.dataset) {
			if c.verbose {
				fmt.Println("All rows have been selected. ")
			}
			break
		}

		if selectedRows[index] {
			if c.verbose {
				fmt.Println("Row already selected. ", selectedRows[index])
			}

			continue
		}

		for class, descriptionList := range c.prompts {
			// if label already has enough prompts, skip to the next label
			if len(descriptionList) >= maxDescriptions {
				if c.verbose {
					fmt.Println("Label already has enough prompts. ", class, len(descriptionList), maxDescriptions)
				}
				continue
			}

			classificationProfile, err := c.GenerateClassifierProfile(class, row, ClassifierProfile{})

			if c.verbose {
				fmt.Println("Classification Profile: ", classificationProfile)
			}

			if err != nil {
				log.Fatal("Train: Failed to generate classifier profile", err)
			} else {
				c.prompts[class] = append(descriptionList, classificationProfile.Description...)
			}
		}

		selectedRows[index] = true
	}

	if c.verbose {
		fmt.Println("Prompts After: ", c.prompts)
	}

	return nil
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
	// convert c.prompts to a string
	classDescriptors := "Class->Description\n"
	for className, descriptionList := range c.prompts {
		for _, description := range descriptionList {
			classDescriptors += fmt.Sprintf("%s: %s\n", className, description)
		}
	}

	// TODO: Implement OpenAI API call
	systemPrompt := fmt.Sprintf(`You are an AI assistant that performs classification. 
	You will be given a map of predicted classes and their corresponding descriptions. 
	Use this information to classify the given data point.
	Respond in JSON with { predicted_class: <class>, "probability": <probability> }. 
	The label should be only from the given labels.
	Context: %s\n`, classDescriptors)

	if text == "" {
		return ClassificationResult{Label: "", Probability: -1}, fmt.Errorf("text cannot be empty")
	}

	userPrompt := fmt.Sprintf(`Classify the following text: "%s"`, text)

	text, err := c.ai.GenerateText(userPrompt, GenerateTextOptions{Verbose: false, System: systemPrompt})

	if err != nil {
		return ClassificationResult{Label: "", Probability: -1}, err
	}

	if c.verbose {
		fmt.Println("System Prompt: ", systemPrompt)
		fmt.Println("Prompt: ", userPrompt)
		fmt.Println("Generated Text: ", text)
	}

	result, err := CleanGPTJson[ClassificationResult](text)

	if err != nil {
		return ClassificationResult{Label: "", Probability: -1}, err
	}

	if c.targetColumn != "" {
		result.Label = c.targetColumn
	} else {
		result.Label = ""
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
