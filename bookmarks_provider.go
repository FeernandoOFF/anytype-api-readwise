package main

// BookmarksProvider is an interface for fetching books and highlights
type BookmarksProvider interface {
	// GetBooks returns a list of books
	GetBooks() ([]ReadwiseBook, error)
	
	// GetHighlights returns a list of highlights for a specific book
	GetHighlights(bookID int) ([]Highlight, error)
}