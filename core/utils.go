package core

import (
	"encoding/json"
	"fmt"
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
