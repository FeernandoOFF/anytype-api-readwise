// config.go
package core

import (
	"fmt"
	"os"
)

type Config struct {
	ReadwiseToken     string
	AnytypeAPIKey     string
	AnytypeBaseURL    string
	AnytypeVersion    string
	TemplatePath      string
	AnytypeTemplateID string
	ObjectType        string
	SpaceID           string
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ValidateConfig(config *Config) error {
	if config.ReadwiseToken == "" {
		return fmt.Errorf("READWISE_TOKEN environment variable is required")
	}
	if config.AnytypeAPIKey == "" {
		return fmt.Errorf("ANYTYPE_API_KEY environment variable is required")
	}

	// Ensure at least one template option is provided
	if config.AnytypeTemplateID == "" {
		// If no Anytype template ID is provided, check for a valid markdown template
		if config.TemplatePath == "" {
			return fmt.Errorf("either a markdown template path or an Anytype template ID must be provided")
		}
		if _, err := os.Stat(config.TemplatePath); os.IsNotExist(err) {
			return fmt.Errorf("template file not found: %s", config.TemplatePath)
		}
	}

	return nil
}
