package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func LogJSON(label string, v interface{}) {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("%s: (error marshaling: %v)", label, err)
		return
	}
	fmt.Printf("\n%s:\n%s\n\n", label, string(jsonBytes))
}
