package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

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
