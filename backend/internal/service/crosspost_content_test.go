package service

import (
	"strings"
	"testing"
)

func TestMarkdownBodyStripsFrontmatter(t *testing.T) {
	raw := "---\ntitle: Hello\nslug: hello\n---\n\n# Hello\n\nBody text.\n"
	got := MarkdownBody(raw)
	if strings.Contains(got, "title:") {
		t.Fatalf("frontmatter leaked: %q", got)
	}
	if !strings.Contains(got, "Body text.") {
		t.Fatalf("missing body: %q", got)
	}
}

func TestMediumTagsLimitsAndDedupes(t *testing.T) {
	got := MediumTags([]string{"Go", "go", "engineering", "personal", "too-long-tag-that-exceeds-limit-here"})
	if len(got) != 3 {
		t.Fatalf("got %v, want 3 tags", got)
	}
	if got[0] != "Go" || got[1] != "engineering" || got[2] != "personal" {
		t.Fatalf("got %v", got)
	}
}

func TestMediumContentAddsTitleCanonicalAndCover(t *testing.T) {
	md := MediumContent("https://shashwatdixit.com", CanonicalPost{
		Slug:       "hello",
		Title:      "Hello",
		CoverImage: "/images/cover.jpg",
		ContentMD:  "---\ntitle: Hello\n---\n\nParagraph.\n",
	})
	for _, want := range []string{
		"![cover](https://shashwatdixit.com/images/cover.jpg)",
		"# Hello",
		"Paragraph.",
		"https://shashwatdixit.com/blog/hello",
		"Originally published",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("missing %q in:\n%s", want, md)
		}
	}
}

func TestAbsolutizeRewritesRootRelativeURLs(t *testing.T) {
	md := AbsolutizeMarkdown("https://example.com", "See [pic](/img.png)")
	if md != "See [pic](https://example.com/img.png)" {
		t.Fatalf("markdown: %q", md)
	}
	html := AbsolutizeHTML("https://example.com", `<img src="/img.png"><a href="/blog">x</a>`)
	if !strings.Contains(html, `src="https://example.com/img.png"`) {
		t.Fatalf("html src: %q", html)
	}
	if !strings.Contains(html, `href="https://example.com/blog"`) {
		t.Fatalf("html href: %q", html)
	}
}

func TestSubstackHTMLIncludesCanonicalFooter(t *testing.T) {
	html := SubstackHTML("https://shashwatdixit.com", CanonicalPost{
		Slug:        "hello",
		Title:       "Hello",
		ContentHTML: "<p>Hi</p>",
		CoverImage:  "https://cdn.example/cover.jpg",
	})
	for _, want := range []string{
		`<img src="https://cdn.example/cover.jpg"`,
		"<p>Hi</p>",
		`href="https://shashwatdixit.com/blog/hello"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q in %s", want, html)
		}
	}
}
