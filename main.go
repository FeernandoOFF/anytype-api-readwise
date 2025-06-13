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
	objectType := flag.String("type", "bookmark", "Anytype object type to create")
	spaceID := flag.String("space", "", "Anytype space ID (optional)")
	flag.Parse()

	// Initialize configuration
	config := &Config{
		ReadwiseToken:  os.Getenv("READWISE_TOKEN"),
		AnytypeAPIKey:  os.Getenv("ANYTYPE_API_KEY"),
		AnytypeBaseURL: getEnvOrDefault("ANYTYPE_API_BASE_URL", "http://localhost:31009"),
		AnytypeVersion: getEnvOrDefault("ANYTYPE_VERSION", "2025-05-20"),
		TemplatePath:   *templatePath,
		ObjectType:     *objectType,
		SpaceID:        *spaceID,
	}

	if err := validateConfig(config); err != nil {
		log.Fatal("Configuration error:", err)
	}

	// Initialize services
	readwiseClient := NewReadwiseClient(config.ReadwiseToken)
	anytypeClient := NewAnytypeClient(config.AnytypeAPIKey, config.AnytypeBaseURL, config.AnytypeVersion)

	// Create syncer and run
	syncer := NewSyncer(readwiseClient, anytypeClient, config)
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
	if _, err := os.Stat(config.TemplatePath); os.IsNotExist(err) {
		return fmt.Errorf("template file not found: %s", config.TemplatePath)
	}
	return nil
}
