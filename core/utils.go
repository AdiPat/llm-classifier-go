package core

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
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
