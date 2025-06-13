package sync

import (
	"anytype-readwise/core"
	"anytype-readwise/feature/bookmarks"
	"anytype-readwise/feature/notes"
	"anytype-readwise/feature/templates"
	"fmt"
	"time"
)

type Syncer struct {
	bookmarksProvider bookmarks.BookmarksProvider
	anytypeClient     *notes.AnytypeClient
	templateProvider  templates.TemplateProvider
	config            *core.Config
}

func NewSyncer(bookmarksProvider bookmarks.BookmarksProvider, anytypeClient *notes.AnytypeClient, templateProvider templates.TemplateProvider, config *core.Config) *Syncer {
	return &Syncer{
		bookmarksProvider: bookmarksProvider,
		anytypeClient:     anytypeClient,
		templateProvider:  templateProvider,
		config:            config,
	}
}

func (s *Syncer) Sync() error {
	fmt.Println("Starting bookmark sync to Anytype...")

	// Get space ID
	spaceID, err := s.anytypeClient.GetSpaceID()
	fmt.Println("Space ID:", spaceID)
	if err != nil {
		return fmt.Errorf("failed to get space ID: %w", err)
	}

	// Fetch books from the bookmarks provider
	fmt.Println("Fetching books from bookmarks provider...")
	books, err := s.bookmarksProvider.GetBooks()
	if err != nil {
		return fmt.Errorf("failed to fetch books: %w", err)
	}

	fmt.Printf("Found %d books\n", len(books))

	//for i, book := range books {
	i := 6
	book := books[i]
	fmt.Printf("Processing book %d/%d: %s by %s\n", i+1, len(books), book.Title, book.Author)

	// Fetch highlights for this book
	highlights, err := s.bookmarksProvider.GetHighlights(book.ID)
	if err != nil {
		fmt.Printf("Warning: failed to fetch highlights for book %s: %v\n", book.Title, err)
		highlights = []bookmarks.Highlight{} // Continue with empty highlights
	}

	fmt.Printf("Found %d highlights\n", len(highlights))

	// Render the template
	templateData := templates.TemplateData{
		Book:       book,
		Highlights: highlights,
		SyncDate:   time.Now().Format("January 2, 2006"),
	}
	content, err := s.templateProvider.Render(templateData)
	if err != nil {
		return fmt.Errorf("failed to render template for book %s: %w", book.Title, err)
	}

	// Create or update object in Anytype
	obj, err := s.anytypeClient.CreateOrUpdateNoteFromBook(spaceID, book, content)
	if err != nil {
		return fmt.Errorf("failed to create or update object for book %s: %w", book.Title, err)
	}

	fmt.Printf("Created or updated object: %s (ID: %s)\n", obj.Name, obj.ID)
	//}

	return nil
}
