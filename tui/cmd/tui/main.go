package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/ssh"
	"golang.org/x/term"

	"github.com/shashwat-dixit/portfolio/tui/internal/api"
	"github.com/shashwat-dixit/portfolio/tui/internal/profile"
	"github.com/shashwat-dixit/portfolio/tui/internal/server"
	"github.com/shashwat-dixit/portfolio/tui/internal/ui"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	listen := flag.String("listen", env("SSH_LISTEN", ""), "SSH listen address (host:port or port). Empty runs locally when stdin is a TTY.")
	hostKey := flag.String("host-key", env("SSH_HOST_KEY", ".ssh/id_ed25519"), "path to the SSH host key")
	apiURL := flag.String("api", env("API_URL", "http://localhost:8080"), "blog API base URL")
	siteURL := flag.String("site", env("SITE_URL", profile.URL), "public website URL")
	local := flag.Bool("local", false, "run the TUI in the current terminal instead of starting SSH")
	flag.Parse()

	addr := normalizeListen(*listen)
	if *local || (addr == "" && term.IsTerminal(int(os.Stdin.Fd()))) {
		if err := runLocal(*apiURL, *siteURL); err != nil {
			slog.Error("tui failed", "error", err)
			os.Exit(1)
		}
		return
	}
	if addr == "" {
		addr = ":2222"
	}

	if err := os.MkdirAll(filepath.Dir(*hostKey), 0o700); err != nil {
		slog.Error("create host key dir", "error", err)
		os.Exit(1)
	}

	srv, err := server.New(server.Config{
		Listen:      addr,
		HostKeyPath: *hostKey,
		APIURL:      *apiURL,
		SiteURL:     *siteURL,
	})
	if err != nil {
		slog.Error("create ssh server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("ssh tui listening", "addr", addr, "api", *apiURL)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			slog.Error("ssh server failed", "error", err)
			done <- syscall.SIGTERM
		}
	}()

	<-done
	slog.Info("stopping ssh tui")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		slog.Error("ssh shutdown", "error", err)
	}
}

func runLocal(apiURL, siteURL string) error {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24
	}
	section, slug := ui.ParseArgs(flag.Args())
	m := ui.New(ui.Options{
		Width:        width,
		Height:       height,
		StartSection: section,
		StartSlug:    slug,
		Client:       api.New(apiURL),
		SiteURL:      siteURL,
	})
	_, err = tea.NewProgram(m).Run()
	return err
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func normalizeListen(v string) string {
	if v == "" {
		v = os.Getenv("SSH_PORT")
	}
	if v == "" {
		return ""
	}
	if strings.HasPrefix(v, ":") || strings.Contains(v, ":") {
		return v
	}
	return ":" + v
}
