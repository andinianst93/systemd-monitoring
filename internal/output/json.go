package output

import (
	"encoding/json"
	"fmt"

	"github.com/andinianst93/systemd-monitoring/internal/models"
)

// PrintJSON prints services in JSON format
func PrintJSON(serviceList *models.ServiceList) error {
	// 1. Marshal to JSON with indentation (pretty print)
	// json.MarshalIndent(data, prefix, indent)
	// prefix = "" (no prefix for each line)
	// indent = "  " (2 spaces for indentation)
	data, err := json.MarshalIndent(serviceList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 2. Print the JSON string
	fmt.Println(string(data))

	return nil
}

// PrintJSONPretty prints a single service in detailed JSON
func PrintJSONPretty(service *models.ServiceInfo) error {
	// Marshal single service to JSON with indentation
	data, err := json.MarshalIndent(service, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Print the JSON string
	fmt.Println(string(data))

	return nil
}
