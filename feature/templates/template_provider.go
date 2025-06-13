package templates

import (
	"anytype-readwise/feature/bookmarks"
)

// TemplateData contains the data needed to render a template
type TemplateData struct {
	Book       bookmarks.ReadwiseBook
	Highlights []bookmarks.Highlight
	SyncDate   string
}

// TemplateProvider is an interface for rendering templates
type TemplateProvider interface {
	// Render renders a template with the given data
	Render(data TemplateData) (string, error)
}
