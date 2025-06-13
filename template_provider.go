package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// TemplateData contains the data needed to render a template
type TemplateData struct {
	Book       ReadwiseBook
	Highlights []Highlight
	SyncDate   string
}

// TemplateProvider is an interface for rendering templates
type TemplateProvider interface {
	// Render renders a template with the given data
	Render(data TemplateData) (string, error)
}

// MarkdownTemplateProvider implements TemplateProvider using a markdown template file
type MarkdownTemplateProvider struct {
	templatePath string
}

// NewMarkdownTemplateProvider creates a new MarkdownTemplateProvider
func NewMarkdownTemplateProvider(templatePath string) *MarkdownTemplateProvider {
	return &MarkdownTemplateProvider{
		templatePath: templatePath,
	}
}

// Render renders the markdown template with the given data
func (p *MarkdownTemplateProvider) Render(data TemplateData) (string, error) {
	// Load the template file
	templateContent, err := os.ReadFile(p.templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	// Define custom functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	// Parse the template
	tmpl, err := template.New("book").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// AnytypeTemplateProvider implements TemplateProvider using an Anytype template
type AnytypeTemplateProvider struct {
	anytypeClient *AnytypeClient
	templateID    string
}

// NewAnytypeTemplateProvider creates a new AnytypeTemplateProvider
func NewAnytypeTemplateProvider(anytypeClient *AnytypeClient, templateID string) *AnytypeTemplateProvider {
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
