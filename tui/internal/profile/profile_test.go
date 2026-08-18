package profile

import (
	"strings"
	"testing"
)

func TestMarkdownIncludesCorePortfolio(t *testing.T) {
	md := Markdown()
	want := []string{
		"# Shashwat Dixit",
		"Bengaluru, India",
		"ssh shashwatdixit.com -p 2222",
		"Interview Kickstart",
		"Instahyre",
		"Pummyz Foods",
		"Nitte Meenakshi Institute of Technology",
		"PostgreSQL",
		"Jamin",
		"Zort",
		"shashwatmain@gmail.com",
		"https://github.com/shashwat-dixit",
	}
	for _, s := range want {
		if !strings.Contains(md, s) {
			t.Errorf("markdown missing %q", s)
		}
	}
}

func TestWorkHasPresentRole(t *testing.T) {
	if Work[0].End != "Present" {
		t.Fatalf("current role end = %q, want Present", Work[0].End)
	}
	if Work[0].Company != "Interview Kickstart" {
		t.Fatalf("current company = %q", Work[0].Company)
	}
}

func TestSkillsAndProjectsNotEmpty(t *testing.T) {
	if len(Skills) < 10 {
		t.Fatalf("skills = %d, want at least 10", len(Skills))
	}
	if len(Projects) == 0 {
		t.Fatal("expected projects")
	}
}
