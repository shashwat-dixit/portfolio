package service

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

type MarkdownService struct {
	md goldmark.Markdown
}

func NewMarkdown() *MarkdownService {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	return &MarkdownService{md: md}
}

func (s *MarkdownService) Parse(raw []byte) (*model.Frontmatter, string, int, error) {
	fm, body, err := splitFrontmatter(raw)
	if err != nil {
		return nil, "", 0, err
	}

	var frontmatter model.Frontmatter
	if err := yaml.Unmarshal(fm, &frontmatter); err != nil {
		return nil, "", 0, fmt.Errorf("parse frontmatter yaml: %w", err)
	}

	var buf bytes.Buffer
	if err := s.md.Convert(body, &buf); err != nil {
		return nil, "", 0, fmt.Errorf("render markdown: %w", err)
	}

	wordCount := countWords(body)

	return &frontmatter, buf.String(), wordCount, nil
}

func splitFrontmatter(raw []byte) ([]byte, []byte, error) {
	content := string(raw)
	content = strings.TrimLeftFunc(content, unicode.IsSpace)

	if !strings.HasPrefix(content, "---") {
		return nil, raw, nil
	}

	content = content[3:]
	end := strings.Index(content, "\n---")
	if end == -1 {
		return nil, raw, fmt.Errorf("unterminated frontmatter")
	}

	fm := []byte(content[:end])
	body := []byte(content[end+4:])

	return fm, body, nil
}

func countWords(text []byte) int {
	words := strings.Fields(string(text))
	return len(words)
}
