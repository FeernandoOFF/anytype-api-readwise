package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Syncer struct {
	readwiseClient *ReadwiseClient
	anytypeClient  *AnytypeClient
	config         *Config
	template       *template.Template
}

type BookTemplateData struct {
	Book       ReadwiseBook
	Highlights []Highlight
	SyncDate   string
}

func NewSyncer(readwiseClient *ReadwiseClient, anytypeClient *AnytypeClient, config *Config) *Syncer {
	return &Syncer{
		readwiseClient: readwiseClient,
		anytypeClient:  anytypeClient,
		config:         config,
	}
}

func (s *Syncer) loadTemplate() error {
	templateContent, err := os.ReadFile(s.config.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Define custom functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	tmpl, err := template.New("book").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	s.template = tmpl
	return nil
}

func (s *Syncer) getSpaceID() (string, error) {
	if s.config.SpaceID != "" {
		return s.config.SpaceID, nil
	}

	spaces, err := s.anytypeClient.GetSpaces()
	if err != nil {
		return "", err
	}

	if len(spaces) == 0 {
		return "", fmt.Errorf("no spaces found")
	}

	defaultSpace := spaces[1]
	// Use the first space if none specified
	fmt.Printf("Using space: %s (ID: %s)\n", defaultSpace.Name, defaultSpace.ID)
	return defaultSpace.ID, nil
}

func (s *Syncer) renderTemplate(book ReadwiseBook, highlights []Highlight) (string, error) {
	if s.template == nil {
		if err := s.loadTemplate(); err != nil {
			return "", err
		}
	}

	data := BookTemplateData{
		Book:       book,
		Highlights: highlights,
		SyncDate:   time.Now().Format("January 2, 2006"),
	}

	var buf strings.Builder
	if err := s.template.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (s *Syncer) Sync() error {
	fmt.Println("Starting Readwise to Anytype sync...")

	// Get space ID
	spaceID, err := s.getSpaceID()
	fmt.Println("Space ID:", spaceID)
	if err != nil {
		return fmt.Errorf("failed to get space ID: %w", err)
	}

	// Fetch books from Readwise
	fmt.Println("Fetching books from Readwise...")
	books, err := s.readwiseClient.GetBooks()
	if err != nil {
		return fmt.Errorf("failed to fetch books: %w", err)
	}

	fmt.Printf("Found %d books\n", len(books))

	//for i, book := range books {
	i := 5
	book := books[i]
	fmt.Printf("Processing book %d/%d: %s by %s\n", i+1, len(books), book.Title, book.Author)

	// Fetch highlights for this book
	highlights, err := s.readwiseClient.GetHighlights(book.ID)
	if err != nil {
		fmt.Printf("Warning: failed to fetch highlights for book %s: %v\n", book.Title, err)
		highlights = []Highlight{} // Continue with empty highlights
	}

	fmt.Printf("Found %d highlights\n", len(highlights))

	//// Render the template
	content, err := s.renderTemplate(book, highlights)
	if err != nil {
		return fmt.Errorf("failed to render template for book %s: %w", book.Title, err)
	}

	//// Create object in Anytype
	createReq := CreateObjectRequest{
		Name:    fmt.Sprintf("%s - %s [SYNC]", book.Title, book.Author),
		TypeKey: s.config.ObjectType,
		Body:    content,
		Icon: &ObjectIcon{
			Emoji:  "📚",
			Format: "emoji",
		},
		Properties: []Property{
			{
				Key:   "description",
				Value: strconv.Itoa(book.ID),
			},
		},
	}
	obj, err := s.anytypeClient.CreateObject(spaceID, createReq)
	if err != nil {
		return fmt.Errorf("failed to create object for book %s: %w", book.Title, err)
	}

	fmt.Printf("Created object: %s (ID: %s)\n", obj.Name, obj.ID)
	//}

	return nil
}
