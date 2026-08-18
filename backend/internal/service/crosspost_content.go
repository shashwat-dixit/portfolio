package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	mdRootRel   = regexp.MustCompile(`\]\((/[^)]+)\)`)
	htmlRootRel = regexp.MustCompile(`(?i)((?:src|href)=["'])(/[^"']+)`)
)

func CanonicalPostURL(siteURL, slug string) string {
	return strings.TrimRight(siteURL, "/") + "/blog/" + slug
}

func AbsoluteURL(siteURL, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	base := strings.TrimRight(siteURL, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}

func AbsolutizeMarkdown(siteURL, md string) string {
	base := strings.TrimRight(siteURL, "/")
	if base == "" {
		return md
	}
	return mdRootRel.ReplaceAllStringFunc(md, func(match string) string {
		parts := mdRootRel.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		return "](" + base + parts[1] + ")"
	})
}

func AbsolutizeHTML(siteURL, html string) string {
	base := strings.TrimRight(siteURL, "/")
	if base == "" {
		return html
	}
	return htmlRootRel.ReplaceAllStringFunc(html, func(match string) string {
		parts := htmlRootRel.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		return parts[1] + base + parts[2]
	})
}

func MediumTags(tags []string) []string {
	out := make([]string, 0, 3)
	seen := make(map[string]bool)
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" || utf8.RuneCountInString(tag) > 25 {
			continue
		}
		key := strings.ToLower(tag)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, tag)
		if len(out) == 3 {
			break
		}
	}
	return out
}

func MediumContent(siteURL string, post CanonicalPost) string {
	body := AbsolutizeMarkdown(siteURL, MarkdownBody(post.ContentMD))
	var b strings.Builder
	if cover := AbsoluteURL(siteURL, post.CoverImage); cover != "" {
		fmt.Fprintf(&b, "![cover](%s)\n\n", cover)
	}
	title := strings.TrimSpace(post.Title)
	if title != "" && !hasLeadingTitle(body, title) {
		fmt.Fprintf(&b, "# %s\n\n", title)
	}
	b.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		b.WriteByte('\n')
	}
	canonical := CanonicalPostURL(siteURL, post.Slug)
	fmt.Fprintf(&b, "\n---\n\n*Originally published at [%s](%s).*\n", hostLabel(siteURL), canonical)
	return b.String()
}

func SubstackHTML(siteURL string, post CanonicalPost) string {
	html := AbsolutizeHTML(siteURL, post.ContentHTML)
	canonical := CanonicalPostURL(siteURL, post.Slug)
	footer := fmt.Sprintf(
		`<hr><p><em>Originally published at <a href="%s">%s</a>.</em></p>`,
		canonical,
		htmlEscape(hostLabel(siteURL)),
	)
	if cover := AbsoluteURL(siteURL, post.CoverImage); cover != "" {
		return fmt.Sprintf(`<p><img src="%s" alt=""></p>%s%s`, htmlEscape(cover), html, footer)
	}
	return html + footer
}

type CanonicalPost struct {
	Slug        string
	Title       string
	ContentMD   string
	ContentHTML string
	CoverImage  string
}

func hasLeadingTitle(body, title string) bool {
	trimmed := strings.TrimSpace(body)
	if !strings.HasPrefix(trimmed, "# ") {
		return false
	}
	firstLine, _, _ := strings.Cut(trimmed, "\n")
	return strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(firstLine, "# ")), title)
}

func hostLabel(siteURL string) string {
	u, err := url.Parse(siteURL)
	if err != nil || u.Host == "" {
		return strings.TrimRight(strings.TrimPrefix(strings.TrimPrefix(siteURL, "https://"), "http://"), "/")
	}
	return u.Host
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&quot;",
	)
	return replacer.Replace(s)
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max])
}
