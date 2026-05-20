package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/service"
)

type FeedHandler struct {
	svc *service.PostService
	cfg *config.Config
}

func NewFeedHandler(svc *service.PostService, cfg *config.Config) *FeedHandler {
	return &FeedHandler{svc: svc, cfg: cfg}
}

func (h *FeedHandler) RSS(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.List(r.Context(), "", 1, 50, false)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	feed := rssRoot{
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: rssChannel{
			Title:       "Shashwat Dixit's Blog",
			Link:        h.cfg.SiteURL + "/blog",
			Description: "Thoughts on software development, life, and more.",
			Language:    "en-us",
			LastBuild:   time.Now().UTC().Format(time.RFC1123Z),
			AtomLink: atomLink{
				Href: h.cfg.SiteURL + "/api/feed.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
		},
	}

	for _, post := range resp.Posts {
		item := rssItem{
			Title:       post.Title,
			Link:        fmt.Sprintf("%s/blog/%s", h.cfg.SiteURL, post.Slug),
			Description: post.Description,
			PubDate:     func() string { if post.PublishedAt != nil { return post.PublishedAt.UTC().Format(time.RFC1123Z) }; return "" }(),
			GUID:        fmt.Sprintf("%s/blog/%s", h.cfg.SiteURL, post.Slug),
		}
		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=86400")
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(feed)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type rssRoot struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	AtomNS  string     `xml:"xmlns:atom,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Language    string   `xml:"language"`
	LastBuild   string   `xml:"lastBuildDate"`
	AtomLink    atomLink `xml:"atom:link"`
	Items       []rssItem
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	GUID        string   `xml:"guid"`
}
