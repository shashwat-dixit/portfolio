package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"

	"github.com/shashwat-dixit/portfolio/tui/internal/api"
	"github.com/shashwat-dixit/portfolio/tui/internal/profile"
)

func renderMarkdown(md string, width int) string {
	if width < 20 {
		width = 20
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return wrap(md, width)
	}
	out, err := r.Render(md)
	if err != nil {
		return wrap(md, width)
	}
	return strings.TrimSpace(out)
}

func wrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	return lipgloss.NewStyle().Width(width).Render(s)
}

func renderHome(width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("Hi, I'm "+strings.Split(profile.Name, " ")[0]))
	fmt.Fprintf(&b, "%s\n\n", wrap(profile.Description, width))
	fmt.Fprintf(&b, "%s\n", dimStyle().Render(profile.Location))
	fmt.Fprintf(&b, "%s\n\n", accentStyle().Render(profile.URL))
	fmt.Fprintf(&b, "%s\n", mutedStyle().Render("This is the same portfolio as the website, in your terminal."))
	fmt.Fprintf(&b, "%s\n", mutedStyle().Render("Use the sidebar or number keys to look around. Blog posts load live from the API."))
	return b.String()
}

func renderAbout(width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("About"))
	b.WriteString(renderMarkdown(profile.Summary, width))
	return b.String()
}

func renderWork(width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("Work Experience"))
	fmt.Fprintf(&b, "%s\n\n", mutedStyle().Render("My professional journey building software at scale."))
	for i, job := range profile.Work {
		if i > 0 {
			b.WriteString("\n")
		}
		badge := ""
		if job.Badge != "" {
			badge = "  " + dimStyle().Render("["+job.Badge+"]")
		}
		fmt.Fprintf(&b, "%s%s\n", accentStyle().Bold(true).Render(job.Company), badge)
		fmt.Fprintf(&b, "%s\n", titleStyle().Render(job.Title))
		fmt.Fprintf(&b, "%s\n", dimStyle().Render(fmt.Sprintf("%s · %s – %s", job.Location, job.Start, job.End)))
		fmt.Fprintf(&b, "%s\n", mutedStyle().Render(job.URL))
		for _, bullet := range job.Bullets {
			fmt.Fprintf(&b, "  %s %s\n", accentStyle().Render("•"), wrap(bullet, max(10, width-4)))
		}
	}
	return b.String()
}

func renderEducation(width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("Education"))
	for _, school := range profile.Education {
		fmt.Fprintf(&b, "%s\n", accentStyle().Bold(true).Render(school.Name))
		fmt.Fprintf(&b, "%s\n", wrap(school.Degree, width))
		fmt.Fprintf(&b, "%s\n", dimStyle().Render(fmt.Sprintf("%s – %s", school.Start, school.End)))
		fmt.Fprintf(&b, "%s\n", mutedStyle().Render(school.URL))
	}
	return b.String()
}

func renderSkills(width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("Skills"))
	row := strings.Builder{}
	rowWidth := 0
	for _, skill := range profile.Skills {
		chip := chipStyle().Render(skill)
		w := lipgloss.Width(chip) + 1
		if rowWidth+w > width && rowWidth > 0 {
			b.WriteString(strings.TrimRight(row.String(), " "))
			b.WriteByte('\n')
			row.Reset()
			rowWidth = 0
		}
		row.WriteString(chip)
		row.WriteByte(' ')
		rowWidth += w
	}
	if row.Len() > 0 {
		b.WriteString(strings.TrimRight(row.String(), " "))
		b.WriteByte('\n')
	}
	return b.String()
}

func renderProjects(width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("Projects"))
	fmt.Fprintf(&b, "%s\n\n", mutedStyle().Render("Tools ranging from AI chat platforms to real-time collaborative whiteboards."))
	for i, project := range profile.Projects {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "%s  %s\n", accentStyle().Bold(true).Render(fmt.Sprintf("%02d  %s", i+1, project.Title)), dimStyle().Render(project.Dates))
		fmt.Fprintf(&b, "%s\n", wrap(project.Description, width))
		if len(project.Technologies) > 0 {
			fmt.Fprintf(&b, "%s\n", dimStyle().Render(strings.Join(project.Technologies, " · ")))
		}
		fmt.Fprintf(&b, "%s\n", mutedStyle().Render(project.URL))
		if project.SourceURL != "" {
			fmt.Fprintf(&b, "%s\n", mutedStyle().Render("source  "+project.SourceURL))
		}
	}
	return b.String()
}

func renderContact(width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("Get in Touch"))
	fmt.Fprintf(&b, "%s\n\n", wrap("Have something in mind? Shoot me an email, connect on LinkedIn, or book a quick call.", width))
	fmt.Fprintf(&b, "%s  %s\n", dimStyle().Render("email   "), accentStyle().Render(profile.Email))
	for _, link := range profile.Social {
		if strings.EqualFold(link.Name, "Email") {
			continue
		}
		label := fmt.Sprintf("%-8s", strings.ToLower(link.Name))
		fmt.Fprintf(&b, "%s  %s\n", dimStyle().Render(label), accentStyle().Render(link.URL))
	}
	fmt.Fprintf(&b, "\n%s\n", mutedStyle().Render("Web    "+profile.URL))
	fmt.Fprintf(&b, "%s\n", mutedStyle().Render("SSH    "+profile.SSHCommand))
	return b.String()
}

func renderBlogList(posts []api.PostSummary, selected int, query string, loading bool, err error, width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", titleStyle().Render("Blog"))
	if query != "" {
		fmt.Fprintf(&b, "%s %s\n\n", dimStyle().Render("filter"), accentStyle().Render(query))
	}
	if loading {
		fmt.Fprintf(&b, "%s\n", mutedStyle().Render("Loading posts…"))
		return b.String()
	}
	if err != nil {
		fmt.Fprintf(&b, "%s\n%s\n", errorStyle().Render("Could not load posts."), mutedStyle().Render(err.Error()))
		return b.String()
	}
	if len(posts) == 0 {
		if query != "" {
			fmt.Fprintf(&b, "%s\n", mutedStyle().Render("No posts match that search."))
		} else {
			fmt.Fprintf(&b, "%s\n", mutedStyle().Render("No published posts yet."))
		}
		return b.String()
	}
	for i, post := range posts {
		marker := "  "
		style := lipgloss.NewStyle()
		if i == selected {
			marker = accentStyle().Render("▸ ")
			style = style.Bold(true)
		}
		date := formatDate(post.Date)
		meta := strings.TrimSpace(strings.Join([]string{date, fmt.Sprintf("%d min", post.ReadingTime)}, " · "))
		fmt.Fprintf(&b, "%s%s\n", marker, style.Render(post.Title))
		if post.Description != "" {
			fmt.Fprintf(&b, "  %s\n", wrap(post.Description, max(10, width-2)))
		}
		tags := strings.Join(post.Tags, ", ")
		if tags != "" {
			meta = meta + " · " + tags
		}
		fmt.Fprintf(&b, "  %s\n\n", dimStyle().Render(meta))
	}
	return b.String()
}

func renderBlogPost(post *api.Post, width int) string {
	if post == nil {
		return mutedStyle().Render("Loading post…")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", titleStyle().Render(post.Title))
	meta := []string{}
	if d := formatDate(post.Date); d != "" {
		meta = append(meta, d)
	}
	if post.ReadingTime > 0 {
		meta = append(meta, fmt.Sprintf("%d min read", post.ReadingTime))
	}
	if len(post.Tags) > 0 {
		meta = append(meta, strings.Join(post.Tags, ", "))
	}
	if len(meta) > 0 {
		fmt.Fprintf(&b, "%s\n", dimStyle().Render(strings.Join(meta, " · ")))
	}
	b.WriteByte('\n')
	body := api.StripFrontmatter(post.ContentMD)
	if strings.TrimSpace(body) == "" {
		body = post.Description
	}
	b.WriteString(renderMarkdown(body, width))
	return b.String()
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("Jan 2, 2006")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
