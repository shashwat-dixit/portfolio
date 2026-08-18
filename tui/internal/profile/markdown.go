package profile

import (
	"fmt"
	"strings"
)

// Markdown returns a plain-text / markdown dump of the portfolio.
// Used when an SSH client has no TTY (ssh -T) so the same content
// still lands in the terminal.
func Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", Name)
	fmt.Fprintf(&b, "> %s\n\n", Description)
	fmt.Fprintf(&b, "%s. [Website](%s) · [Email](mailto:%s)\n\n", Location, URL, Email)
	fmt.Fprintf(&b, "Terminal UI: `%s`\n\n", SSHCommand)

	fmt.Fprintf(&b, "## About\n\n%s\n\n", Summary)

	fmt.Fprintf(&b, "## Work\n\n")
	for _, job := range Work {
		end := job.End
		if end == "" {
			end = "Present"
		}
		fmt.Fprintf(&b, "### [%s](%s) — %s\n\n", job.Company, job.URL, job.Title)
		fmt.Fprintf(&b, "%s. %s – %s\n\n", job.Location, job.Start, end)
		for _, bullet := range job.Bullets {
			fmt.Fprintf(&b, "- %s\n", bullet)
		}
		b.WriteByte('\n')
	}

	fmt.Fprintf(&b, "## Education\n\n")
	for _, school := range Education {
		fmt.Fprintf(&b, "### [%s](%s)\n\n", school.Name, school.URL)
		fmt.Fprintf(&b, "%s. %s – %s\n\n", school.Degree, school.Start, school.End)
	}

	fmt.Fprintf(&b, "## Skills\n\n%s\n\n", strings.Join(Skills, ", "))

	fmt.Fprintf(&b, "## Projects\n\n")
	for _, project := range Projects {
		fmt.Fprintf(&b, "### [%s](%s)\n\n", project.Title, project.URL)
		fmt.Fprintf(&b, "%s\n", project.Description)
		if len(project.Technologies) > 0 {
			fmt.Fprintf(&b, "\nTech: %s\n", strings.Join(project.Technologies, ", "))
		}
		if project.SourceURL != "" {
			fmt.Fprintf(&b, "\n- [Source](%s)\n", project.SourceURL)
		}
		b.WriteByte('\n')
	}

	fmt.Fprintf(&b, "## Contact\n\n")
	fmt.Fprintf(&b, "- Email: %s\n", Email)
	for _, link := range Social {
		if strings.EqualFold(link.Name, "Email") {
			continue
		}
		fmt.Fprintf(&b, "- [%s](%s)\n", link.Name, link.URL)
	}
	b.WriteByte('\n')
	return b.String()
}
