package util

import (
	"encoding/json"
)

// PackJSON packs a map of strings as a JSON string.
func PackJSON(mapped map[string]string) string {
	jsonData, _ := json.Marshal(mapped)
	return string(jsonData)
}
