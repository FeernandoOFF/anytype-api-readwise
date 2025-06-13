package templates

import (
	"anytype-readwise/feature/notes"
	"fmt"
)

// AnytypeTemplateProvider implements TemplateProvider using an Anytype template
type AnytypeTemplateProvider struct {
	anytypeClient *notes.AnytypeClient
	templateID    string
}

// NewAnytypeTemplateProvider creates a new AnytypeTemplateProvider
func NewAnytypeTemplateProvider(anytypeClient *notes.AnytypeClient, templateID string) *AnytypeTemplateProvider {
	return &AnytypeTemplateProvider{
		anytypeClient: anytypeClient,
		templateID:    templateID,
	}
}

// Render renders the Anytype template with the given data
func (p *AnytypeTemplateProvider) Render(data TemplateData) (string, error) {
	// For now, this is a placeholder implementation
	// In a real implementation, this would fetch a template from Anytype using the template ID
	// and apply the data to it

	// Create a simple markdown template as a fallback
	templateStr := fmt.Sprintf("# %s\n\n", data.Book.Title)
	templateStr += fmt.Sprintf("**Author:** %s\n\n", data.Book.Author)
	templateStr += fmt.Sprintf("**Category:** %s\n\n", data.Book.Category)
	templateStr += fmt.Sprintf("**Source:** %s\n\n", data.Book.Source)
	templateStr += fmt.Sprintf("**Synced on:** %s\n\n", data.SyncDate)

	templateStr += "## Highlights\n\n"

	for i, highlight := range data.Highlights {
		templateStr += fmt.Sprintf("### Highlight %d\n\n", i+1)
		templateStr += fmt.Sprintf("%s\n\n", highlight.Text)
		if highlight.Note != "" {
			templateStr += fmt.Sprintf("**Note:** %s\n\n", highlight.Note)
		}
	}

	return templateStr, nil
}
