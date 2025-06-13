// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Command line flags
	templatePath := flag.String("template", "book_template.md", "Path to markdown template file")
	anytypeTemplateID := flag.String("anytype-template", "", "Anytype template ID (optional)")
	objectType := flag.String("type", "bookmark", "Anytype object type to create")
	spaceID := flag.String("space", "", "Anytype space ID (optional)")
	flag.Parse()

	// Initialize configuration
	config := &Config{
		ReadwiseToken:     os.Getenv("READWISE_TOKEN"),
		AnytypeAPIKey:     os.Getenv("ANYTYPE_API_KEY"),
		AnytypeBaseURL:    getEnvOrDefault("ANYTYPE_API_BASE_URL", "http://localhost:31009"),
		AnytypeVersion:    getEnvOrDefault("ANYTYPE_VERSION", "2025-05-20"),
		TemplatePath:      *templatePath,
		AnytypeTemplateID: *anytypeTemplateID,
		ObjectType:        *objectType,
		SpaceID:           *spaceID,
	}

	if err := validateConfig(config); err != nil {
		log.Fatal("Configuration error:", err)
	}

	// Initialize services
	// Create a BookmarksProvider (ReadwiseClient)
	bookmarksProvider := NewReadwiseClient(config.ReadwiseToken)

	// Create an AnytypeClient
	anytypeClient := NewAnytypeClient(config.AnytypeAPIKey, config.AnytypeBaseURL, config.AnytypeVersion, config)

	// Create a TemplateProvider based on configuration
	var templateProvider TemplateProvider
	if config.AnytypeTemplateID != "" {
		// Use AnytypeTemplateProvider if a template ID is provided
		templateProvider = NewAnytypeTemplateProvider(anytypeClient, config.AnytypeTemplateID)
		fmt.Println("Using Anytype template with ID:", config.AnytypeTemplateID)
	} else {
		// Use MarkdownTemplateProvider as fallback
		templateProvider = NewMarkdownTemplateProvider(config.TemplatePath)
		fmt.Println("Using markdown template from:", config.TemplatePath)
	}

	// Create syncer and run
	syncer := NewSyncer(bookmarksProvider, anytypeClient, templateProvider, config)
	if err := syncer.Sync(); err != nil {
		log.Fatal("Sync failed:", err)
	}

	fmt.Println("Sync completed successfully!")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func validateConfig(config *Config) error {
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
