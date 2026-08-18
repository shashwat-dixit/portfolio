package ui

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/shashwat-dixit/portfolio/tui/internal/api"
	"github.com/shashwat-dixit/portfolio/tui/internal/profile"
)

type Section int

const (
	SectionHome Section = iota
	SectionAbout
	SectionWork
	SectionEducation
	SectionSkills
	SectionProjects
	SectionBlog
	SectionContact
	sectionCount
)

var sectionLabels = []string{
	"Home",
	"About",
	"Work",
	"Education",
	"Skills",
	"Projects",
	"Blog",
	"Contact",
}

var sectionKeys = []string{
	"home",
	"about",
	"work",
	"education",
	"skills",
	"projects",
	"blog",
	"contact",
}

type Options struct {
	Width        int
	Height       int
	StartSection Section
	StartSlug    string
	Client       *api.Client
	SiteURL      string
}

type Model struct {
	width, height int
	section       Section
	viewport      viewport.Model
	client        *api.Client
	siteURL       string

	posts        []api.PostSummary
	visiblePosts []api.PostSummary
	selectedPost int
	postsLoading bool
	postsErr     error

	readingPost bool
	post        *api.Post
	postLoading bool
	postErr     error
	startSlug   string

	searching bool
	search    textinput.Model
}

type postsMsg struct {
	posts []api.PostSummary
	err   error
}

type postMsg struct {
	post *api.Post
	err  error
}

func New(opts Options) Model {
	ti := textinput.New()
	ti.Placeholder = "search posts"
	ti.CharLimit = 80
	ti.Prompt = "/ "

	m := Model{
		width:        opts.Width,
		height:       opts.Height,
		section:      opts.StartSection,
		client:       opts.Client,
		siteURL:      strings.TrimRight(opts.SiteURL, "/"),
		startSlug:    opts.StartSlug,
		postsLoading: true,
		search:       ti,
		viewport:     viewport.New(viewport.WithWidth(40), viewport.WithHeight(10)),
	}
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}
	if m.siteURL == "" {
		m.siteURL = profile.URL
	}
	m.viewport.SoftWrap = true
	m.viewport.MouseWheelEnabled = true
	m.syncLayout()
	return m
}

func ParseArgs(args []string) (Section, string) {
	if len(args) == 0 {
		return SectionHome, ""
	}
	raw := strings.TrimSpace(strings.Join(args, " "))
	raw = strings.TrimPrefix(raw, "/")
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == '/' || r == ' '
	})
	if len(parts) == 0 {
		return SectionHome, ""
	}
	name := strings.ToLower(parts[0])
	slug := ""
	if len(parts) > 1 {
		slug = strings.Join(parts[1:], "-")
	}

	if n := sectionNumber(name); n >= 0 {
		sec := Section(n)
		if sec == SectionBlog {
			return sec, slug
		}
		return sec, ""
	}

	for i, key := range sectionKeys {
		if name == key || strings.HasPrefix(key, name) {
			if Section(i) == SectionBlog {
				return SectionBlog, slug
			}
			return Section(i), ""
		}
	}
	// Treat unknown first arg as a blog slug.
	return SectionBlog, name
}

func sectionNumber(s string) int {
	if len(s) != 1 || s[0] < '1' || s[0] > '8' {
		return -1
	}
	return int(s[0] - '1')
}

func (m Model) Init() tea.Cmd {
	return m.fetchPosts()
}

func (m Model) fetchPosts() tea.Cmd {
	if m.client == nil {
		return func() tea.Msg {
			return postsMsg{err: fmt.Errorf("blog API is not configured")}
		}
	}
	client := m.client
	return func() tea.Msg {
		posts, err := client.ListPosts(context.Background())
		return postsMsg{posts: posts, err: err}
	}
}

func (m Model) fetchPost(slug string) tea.Cmd {
	if m.client == nil {
		return func() tea.Msg {
			return postMsg{err: fmt.Errorf("blog API is not configured")}
		}
	}
	client := m.client
	return func() tea.Msg {
		post, err := client.GetPost(context.Background(), slug)
		return postMsg{post: post, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncLayout()
		return m, nil

	case postsMsg:
		m.postsLoading = false
		m.postsErr = msg.err
		m.posts = msg.posts
		slug := m.startSlug
		m.startSlug = ""
		m.applyFilter()
		m.syncLayout()
		if slug != "" {
			return m, m.openNamed(slug)
		}
		return m, nil

	case postMsg:
		m.postLoading = false
		m.postErr = msg.err
		m.post = msg.post
		if msg.err == nil && msg.post != nil {
			m.readingPost = true
		}
		m.syncLayout()
		return m, nil

	case tea.KeyPressMsg:
		return m.updateKey(msg)
	}

	if !m.searching && (m.section != SectionBlog || m.readingPost) {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) updateKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.searching {
		switch key {
		case "esc", "ctrl+c":
			m.searching = false
			m.search.Blur()
			m.search.SetValue("")
			m.applyFilter()
			m.syncLayout()
			return m, nil
		case "enter":
			m.searching = false
			m.search.Blur()
			m.applyFilter()
			m.syncLayout()
			return m, nil
		default:
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			m.applyFilter()
			m.syncLayout()
			return m, cmd
		}
	}

	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "tab":
		m.setSection((m.section + 1) % sectionCount)
		return m, nil
	case "shift+tab":
		m.setSection((m.section - 1 + sectionCount) % sectionCount)
		return m, nil
	case "1", "2", "3", "4", "5", "6", "7", "8":
		m.setSection(Section(sectionNumber(key)))
		return m, nil
	case "esc", "backspace":
		if m.readingPost {
			m.readingPost = false
			m.post = nil
			m.postErr = nil
			m.syncLayout()
			return m, nil
		}
	case "enter":
		if m.section == SectionBlog && !m.readingPost && len(m.visiblePosts) > 0 {
			slug := m.visiblePosts[m.selectedPost].Slug
			return m, m.openNamed(slug)
		}
	case "/":
		if m.section == SectionBlog && !m.readingPost {
			m.searching = true
			m.search.Focus()
			m.syncLayout()
			return m, nil
		}
	case "j", "down":
		if m.section == SectionBlog && !m.readingPost {
			if m.selectedPost < len(m.visiblePosts)-1 {
				m.selectedPost++
				m.syncLayout()
			}
			return m, nil
		}
	case "k", "up":
		if m.section == SectionBlog && !m.readingPost {
			if m.selectedPost > 0 {
				m.selectedPost--
				m.syncLayout()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *Model) setSection(section Section) {
	m.section = section
	m.readingPost = false
	m.post = nil
	m.postErr = nil
	m.searching = false
	m.search.Blur()
	m.syncLayout()
	m.viewport.GotoTop()
}

func (m *Model) applyFilter() {
	q := strings.ToLower(strings.TrimSpace(m.search.Value()))
	if q == "" {
		m.visiblePosts = append([]api.PostSummary(nil), m.posts...)
	} else {
		m.visiblePosts = m.visiblePosts[:0]
		for _, post := range m.posts {
			hay := strings.ToLower(post.Title + " " + post.Description + " " + strings.Join(post.Tags, " "))
			if strings.Contains(hay, q) {
				m.visiblePosts = append(m.visiblePosts, post)
			}
		}
	}
	if m.selectedPost >= len(m.visiblePosts) {
		m.selectedPost = max(0, len(m.visiblePosts)-1)
	}
}

func (m *Model) openNamed(slug string) tea.Cmd {
	if slug == "" {
		return nil
	}
	m.section = SectionBlog
	m.postLoading = true
	m.readingPost = true
	m.post = nil
	m.postErr = nil
	m.syncLayout()
	m.viewport.GotoTop()
	return m.fetchPost(slug)
}

func (m *Model) syncLayout() {
	navW, contentW, contentH := m.paneSizes()
	_ = navW
	m.viewport.SetWidth(contentW)
	m.viewport.SetHeight(contentH)
	m.search.SetWidth(max(10, contentW-2))
	m.viewport.SetContent(m.paneContent(contentW))
}

func (m Model) paneSizes() (navW, contentW, contentH int) {
	navW = 16
	headerH := 4
	footerH := 2
	if m.searching {
		footerH = 3
	}
	contentH = max(3, m.height-headerH-footerH)
	contentW = max(20, m.width-navW-3)
	if m.width < 70 {
		navW = 0
		contentW = max(20, m.width-2)
		contentH = max(3, m.height-headerH-footerH-1)
	}
	return navW, contentW, contentH
}

func (m Model) paneContent(width int) string {
	switch {
	case m.section == SectionBlog && m.readingPost:
		if m.postErr != nil {
			return errorStyle().Render("Could not open that post.") + "\n" + mutedStyle().Render(m.postErr.Error())
		}
		if m.postLoading && m.post == nil {
			return mutedStyle().Render("Loading post…")
		}
		body := renderBlogPost(m.post, width)
		if m.post != nil && m.siteURL != "" {
			body += "\n\n" + mutedStyle().Render("web  "+m.siteURL+"/blog/"+m.post.Slug)
		}
		return body
	case m.section == SectionBlog:
		return renderBlogList(m.visiblePosts, m.selectedPost, m.search.Value(), m.postsLoading, m.postsErr, width)
	case m.section == SectionAbout:
		return renderAbout(width)
	case m.section == SectionWork:
		return renderWork(width)
	case m.section == SectionEducation:
		return renderEducation(width)
	case m.section == SectionSkills:
		return renderSkills(width)
	case m.section == SectionProjects:
		return renderProjects(width)
	case m.section == SectionContact:
		return renderContact(width)
	default:
		return renderHome(width)
	}
}

func (m Model) View() tea.View {
	if m.width < 40 || m.height < 12 {
		v := tea.NewView(mutedStyle().Render("Terminal too small — resize, or use:\n\n  " + profile.SSHCommand + "\n"))
		v.AltScreen = true
		return v
	}

	header := m.renderHeader()
	body := m.renderBody()
	footer := m.renderFooter()
	content := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.WindowTitle = profile.Name
	return v
}

func (m Model) renderHeader() string {
	left := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle().Render(profile.Name),
		dimStyle().Render(profile.Location+" · "+profile.URL),
	)
	right := lipgloss.JoinVertical(lipgloss.Right,
		accentStyle().Render("ssh"),
		dimStyle().Render("portfolio"),
	)
	gap := max(1, m.width-lipgloss.Width(left)-lipgloss.Width(right)-4)
	row := lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", gap), right)
	return headerStyle().Width(m.width - 2).Render(row)
}

func (m Model) renderBody() string {
	navW, contentW, contentH := m.paneSizes()
	pane := m.viewport.View()
	if navW == 0 {
		tabs := make([]string, 0, sectionCount)
		for i, label := range sectionLabels {
			tabs = append(tabs, navItemStyle(Section(i) == m.section).Render(fmt.Sprintf("%d %s", i+1, label)))
		}
		bar := strings.Join(tabs, " ")
		return lipgloss.JoinVertical(lipgloss.Left, bar, pane)
	}
	var nav strings.Builder
	for i, label := range sectionLabels {
		line := fmt.Sprintf("%d  %s", i+1, label)
		nav.WriteString(navItemStyle(Section(i) == m.section).Width(navW - 1).Render(line))
		if i < len(sectionLabels)-1 {
			nav.WriteByte('\n')
		}
	}
	navBox := lipgloss.NewStyle().Width(navW).Height(contentH).Render(nav.String())
	contentBox := lipgloss.NewStyle().Width(contentW).Height(contentH).Render(pane)
	return lipgloss.JoinHorizontal(lipgloss.Top, navBox, contentBox)
}

func (m Model) renderFooter() string {
	if m.searching {
		return footerStyle().Render(m.search.View())
	}
	help := "1-8 jump  tab section  j/k scroll  q quit"
	if m.section == SectionBlog && !m.readingPost {
		help = "j/k select  enter open  / search  tab section  q quit"
	}
	if m.readingPost {
		help = "j/k scroll  esc back  q quit"
	}
	return footerStyle().Width(m.width).Render(help)
}

// CurrentSection is exported for tests.
func (m Model) CurrentSection() Section { return m.section }

func (m Model) VisiblePostCount() int { return len(m.visiblePosts) }

func (m Model) ReadingPost() bool { return m.readingPost }
