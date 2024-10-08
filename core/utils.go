package core

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type RowItem = map[string]string

func CleanGPTJson[T any](jsonStr string) (T, error) {
	var result T

	if jsonStr == "" {
		return result, fmt.Errorf("input JSON string is empty")
	}

	// Regular expression to capture JSON content inside ```json ... ```
	re := regexp.MustCompile("(?s).*?```json\\s*(.*?)\\s*```.*")
	matches := re.FindStringSubmatch(jsonStr)
	if len(matches) >= 2 {
		jsonStr = matches[1]
	}

	// Remove the ```json\n at the beginning and \n``` at the end
	re2 := regexp.MustCompile("^```json\\n|\\n```$")
	jsonStr = re2.ReplaceAllString(jsonStr, "")

	// Remove all newlines and carriage returns
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\r", "")

	// Unmarshal the cleaned JSON string into the provided type
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func ReadCSVFile(filePath string) ([]RowItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var records []RowItem

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		recordMap := RowItem{}
		for i, header := range headers {
			recordMap[header] = record[i]
		}

		records = append(records, recordMap)
	}

	return records, nil
}

func SelectRandomRow(dataset []RowItem) (RowItem, int) {
	if len(dataset) == 0 {
		return nil, -1 // Return nil if the dataset is empty
	}

	randomIndex := rand.Intn(len(dataset))
	return dataset[randomIndex], randomIndex
}

func CountSelectedRows(dataset []RowItem, selected map[int]bool) int {
	count := 0

	for index, _ := range dataset {
		if selected[index] {
			count++
		}
	}
	return count
}

// contains checks if a slice contains a specific element
func Contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func ExtractClasses(dataset []RowItem, targetColumn string) []string {
	classes := []string{}

	for _, row := range dataset {
		class := row[targetColumn]
		if !Contains(classes, class) {
			classes = append(classes, class)
		}
	}

	return classes
}

func CreateFolderIfNotExists(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
