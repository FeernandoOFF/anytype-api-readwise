package main

import (
	"anytype-readwise/core"
	"anytype-readwise/feature/bookmarks"
	"anytype-readwise/feature/notes"
	"anytype-readwise/feature/sync"
	"anytype-readwise/feature/templates"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Command line flags
	templatePath := flag.String("template", "book_template.md", "Path to markdown template file")
	anytypeTemplateID := flag.String("anytype-template", "", "Anytype template ID (optional)")
	objectType := flag.String("type", "Bookmark", "Anytype object type to create")
	spaceID := flag.String("space", "", "Anytype space ID (optional)")
	flag.Parse()

	// Initialize configuration
	config := &core.Config{
		ReadwiseToken:     os.Getenv("READWISE_TOKEN"),
		AnytypeAPIKey:     os.Getenv("ANYTYPE_API_KEY"),
		AnytypeBaseURL:    core.GetEnvOrDefault("ANYTYPE_API_BASE_URL", "http://localhost:31009"),
		AnytypeVersion:    core.GetEnvOrDefault("ANYTYPE_VERSION", "2025-05-20"),
		TemplatePath:      *templatePath,
		AnytypeTemplateID: *anytypeTemplateID,
		ObjectType:        *objectType,
		SpaceID:           *spaceID,
	}

	if err := core.ValidateConfig(config); err != nil {
		log.Fatal("Configuration error:", err)
	}

	// Initialize services
	// Create a BookmarksProvider (ReadwiseClient)
	bookmarksProvider := bookmarks.NewReadwiseClient(config.ReadwiseToken)

	// Create an AnytypeClient
	anytypeClient := notes.NewAnytypeClient(config.AnytypeAPIKey, config.AnytypeBaseURL, config.AnytypeVersion, config)

	// Create a TemplateProvider based on configuration
	var templateProvider templates.TemplateProvider
	if config.AnytypeTemplateID != "" {
		// Use AnytypeTemplateProvider if a template ID is provided
		templateProvider = templates.NewAnytypeTemplateProvider(anytypeClient, config.AnytypeTemplateID)
		fmt.Println("Using Anytype template with ID:", config.AnytypeTemplateID)
	} else {
		// Use MarkdownTemplateProvider as fallback
		templateProvider = templates.NewMarkdownTemplateProvider(config.TemplatePath)
		fmt.Println("Using markdown template from:", config.TemplatePath)
	}

	// Create syncer and run
	syncer := sync.NewSyncer(bookmarksProvider, anytypeClient, templateProvider, config)
	if err := syncer.Sync(); err != nil {
		log.Fatal("Sync failed:", err)
	}

	fmt.Println("Sync completed successfully!")
}
