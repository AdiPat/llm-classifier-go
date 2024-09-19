package core

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func CleanGPTJson[T any](jsonStr string) (T, error) {
	var result T

	if jsonStr == "" {
		return result, fmt.Errorf("input JSON string is empty")
	}

	// Remove the ```json\n at the beginning and \n``` at the end
	re := regexp.MustCompile("^```json\\n|\\n```$")
	jsonStr = re.ReplaceAllString(jsonStr, "")

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

func ReadCSVFile(filePath string) ([]map[string]string, error) {
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

	var records []map[string]string

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		recordMap := make(map[string]string)
		for i, header := range headers {
			recordMap[header] = record[i]
		}

		records = append(records, recordMap)
	}

	return records, nil
}
