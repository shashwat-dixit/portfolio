package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/shashwat-dixit/portfolio/tui/internal/api"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		args    []string
		section Section
		slug    string
	}{
		{nil, SectionHome, ""},
		{[]string{}, SectionHome, ""},
		{[]string{"blog"}, SectionBlog, ""},
		{[]string{"blog", "hello-world"}, SectionBlog, "hello-world"},
		{[]string{"blog/hello-world"}, SectionBlog, "hello-world"},
		{[]string{"work"}, SectionWork, ""},
		{[]string{"3"}, SectionWork, ""},
		{[]string{"about"}, SectionAbout, ""},
		{[]string{"contact"}, SectionContact, ""},
		{[]string{"my-post"}, SectionBlog, "my-post"},
	}
	for _, tt := range tests {
		sec, slug := ParseArgs(tt.args)
		if sec != tt.section || slug != tt.slug {
			t.Errorf("ParseArgs(%q) = %v, %q; want %v, %q", tt.args, sec, slug, tt.section, tt.slug)
		}
	}
}

func TestNewShowsName(t *testing.T) {
	m := New(Options{Width: 100, Height: 30})
	view := m.View()
	if !strings.Contains(view.Content, "Shashwat Dixit") {
		t.Fatalf("view missing name:\n%s", view.Content)
	}
	if !view.AltScreen {
		t.Fatal("expected alt screen")
	}
}

func TestTabAndNumberKeysChangeSection(t *testing.T) {
	m := New(Options{Width: 100, Height: 30})
	updated, _ := m.Update(tea.KeyPressMsg{Code: '2', Text: "2"})
	mod := updated.(Model)
	if mod.CurrentSection() != SectionAbout {
		t.Fatalf("section = %v, want about", mod.CurrentSection())
	}

	updated, _ = mod.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	mod = updated.(Model)
	if mod.CurrentSection() != SectionWork {
		t.Fatalf("section after tab = %v, want work", mod.CurrentSection())
	}
}

func TestBlogFilter(t *testing.T) {
	m := New(Options{Width: 100, Height: 30})
	m.postsLoading = false
	m.posts = []api.PostSummary{
		{Slug: "go", Title: "Writing Go", Tags: []string{"code"}},
		{Slug: "food", Title: "Best dosa", Tags: []string{"food"}},
	}
	m.section = SectionBlog
	m.search.SetValue("dosa")
	m.applyFilter()
	if m.VisiblePostCount() != 1 {
		t.Fatalf("visible = %d, want 1", m.VisiblePostCount())
	}
}

func TestTooSmallTerminal(t *testing.T) {
	m := New(Options{Width: 20, Height: 8})
	view := m.View()
	if !strings.Contains(view.Content, "too small") {
		t.Fatalf("expected too-small message, got %q", view.Content)
	}
}

func TestWorkSectionRendersCompany(t *testing.T) {
	m := New(Options{Width: 100, Height: 40, StartSection: SectionWork})
	view := m.View()
	if !strings.Contains(view.Content, "Interview Kickstart") {
		t.Fatalf("work view missing company:\n%s", view.Content)
	}
}
