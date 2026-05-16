package service

import (
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
)

type MarkdownService struct {
	// TODO: hold a goldmark.Markdown instance configured with:
	// - GFM extension (tables, strikethrough, autolinks)
	// - syntax highlighting via Chroma
	// - heading IDs
}

func NewMarkdown() *MarkdownService {
	// TODO: initialize goldmark with extensions
	return &MarkdownService{}
}

// Parse splits a raw markdown file into frontmatter + body,
// parses the YAML frontmatter, and renders the body to HTML.
//
// Returns the parsed frontmatter, rendered HTML string, and word count.
func (s *MarkdownService) Parse(raw []byte) (*model.Frontmatter, string, int, error) {
	// TODO:
	// 1. Split on "---" delimiters to extract YAML block
	// 2. yaml.Unmarshal into model.Frontmatter
	// 3. Render markdown body to HTML via goldmark
	// 4. Count words in body for reading time calculation
	return nil, "", 0, nil
}
