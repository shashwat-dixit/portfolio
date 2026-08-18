package server

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"

	"github.com/shashwat-dixit/portfolio/tui/internal/api"
	"github.com/shashwat-dixit/portfolio/tui/internal/profile"
	"github.com/shashwat-dixit/portfolio/tui/internal/ui"
)

type Config struct {
	Listen      string
	HostKeyPath string
	APIURL      string
	SiteURL     string
	IdleTimeout time.Duration
	MaxTimeout  time.Duration
}

func New(cfg Config) (*ssh.Server, error) {
	if cfg.Listen == "" {
		cfg.Listen = ":2222"
	}
	if cfg.HostKeyPath == "" {
		cfg.HostKeyPath = ".ssh/id_ed25519"
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 5 * time.Minute
	}
	if cfg.MaxTimeout == 0 {
		cfg.MaxTimeout = 30 * time.Minute
	}

	client := api.New(cfg.APIURL)
	siteURL := cfg.SiteURL

	s, err := wish.NewServer(
		wish.WithAddress(cfg.Listen),
		wish.WithHostKeyPath(cfg.HostKeyPath),
		wish.WithVersion("shashwatdixit-tui"),
		wish.WithIdleTimeout(cfg.IdleTimeout),
		wish.WithMaxTimeout(cfg.MaxTimeout),
		wish.WithPublicKeyAuth(func(ssh.Context, ssh.PublicKey) bool { return true }),
		wish.WithPasswordAuth(func(ssh.Context, string) bool { return true }),
		wish.WithMiddleware(
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				pty, _, _ := sess.Pty()
				section, slug := ui.ParseArgs(sess.Command())
				return ui.New(ui.Options{
					Width:        pty.Window.Width,
					Height:       pty.Window.Height,
					StartSection: section,
					StartSlug:    slug,
					Client:       client,
					SiteURL:      siteURL,
				}), nil
			}),
			plainDumpMiddleware(client),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, err
	}

	s.LocalPortForwardingCallback = func(ssh.Context, string, uint32) bool { return false }
	s.ReversePortForwardingCallback = func(ssh.Context, string, uint32) bool { return false }

	return s, nil
}

// plainDumpMiddleware serves a markdown dump when the client has no TTY
// (for example `ssh -T host`). Interactive sessions fall through to Bubble Tea.
func plainDumpMiddleware(client *api.Client) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			_, _, hasPty := sess.Pty()
			if hasPty {
				next(sess)
				return
			}
			if err := writePlain(sess, client); err != nil {
				slog.Error("plain ssh dump failed", "error", err)
			}
		}
	}
}

func writePlain(sess ssh.Session, client *api.Client) error {
	section, slug := ui.ParseArgs(sess.Command())
	ctx := sess.Context()

	if section != ui.SectionBlog {
		_, err := io.WriteString(sess, profile.Markdown())
		return err
	}

	if slug == "" {
		posts, err := client.ListPosts(ctx)
		if err != nil {
			_, werr := fmt.Fprintf(sess, "Could not load blog posts: %v\n", err)
			return werr
		}
		_, err = io.WriteString(sess, formatBlogIndex(posts))
		return err
	}

	post, err := client.GetPost(ctx, slug)
	if err != nil {
		_, werr := fmt.Fprintf(sess, "Could not load post %q: %v\n", slug, err)
		return werr
	}
	body := strings.TrimSpace(post.ContentMD)
	if body == "" {
		body = "# " + post.Title + "\n\n" + post.Description
	}
	_, err = io.WriteString(sess, body+"\n")
	return err
}

func formatBlogIndex(posts []api.PostSummary) string {
	var b strings.Builder
	b.WriteString("# Blog\n\n")
	if len(posts) == 0 {
		b.WriteString("No published posts yet.\n")
		return b.String()
	}
	for _, post := range posts {
		fmt.Fprintf(&b, "- %s", post.Title)
		if post.Slug != "" {
			fmt.Fprintf(&b, " (%s)", post.Slug)
		}
		b.WriteByte('\n')
		if post.Description != "" {
			fmt.Fprintf(&b, "  %s\n", post.Description)
		}
	}
	return b.String()
}
