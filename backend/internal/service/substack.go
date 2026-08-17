package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
)

type SubstackClient struct {
	publicationURL string
	siteURL        string
	cookie         string
	publish        bool
	baseURL        string
	profileURL     string
	http           httpDoer
	userID         int
	pubUserID      int
}

func NewSubstackClient(cfg *config.Config) *SubstackClient {
	pubURL := strings.TrimRight(cfg.SubstackPublicationURL, "/")
	return &SubstackClient{
		publicationURL: pubURL,
		siteURL:        cfg.SiteURL,
		cookie:         substackCookieHeader(cfg.SubstackCookies, cfg.SubstackSID, cfg.SubstackConnectSID),
		publish:        cfg.SubstackPublish,
		baseURL:        pubURL + "/api/v1",
		profileURL:     "https://substack.com/api/v1/user/profile/self",
		http:           defaultHTTPClient(),
	}
}

func (c *SubstackClient) Enabled() bool {
	return c != nil && c.publicationURL != "" && c.cookie != ""
}

func (c *SubstackClient) SupportsUpdate() bool { return true }

func (c *SubstackClient) Publish(ctx context.Context, post model.SyndicatablePost) (remoteID, remoteURL string, err error) {
	if err := c.ensureProfile(ctx); err != nil {
		return "", "", err
	}
	var created substackDraft
	if err := doJSON(ctx, c.http, http.MethodPost, c.baseURL+"/drafts", c.headers(), c.draftPayload(post), &created); err != nil {
		return "", "", fmt.Errorf("substack create draft: %w", err)
	}
	id := created.id()
	if id == "" {
		return "", "", fmt.Errorf("substack create draft: empty id")
	}
	if c.publish {
		published, pubErr := c.publishDraft(ctx, id)
		if pubErr != nil {
			slog.Warn("substack draft created but publish failed", "id", id, "error", pubErr)
			return id, c.draftURL(id), nil
		}
		if published != "" {
			return id, published, nil
		}
	}
	return id, firstNonEmpty(created.publicURL(c.publicationURL), c.draftURL(id)), nil
}

func (c *SubstackClient) Update(ctx context.Context, remoteID string, post model.SyndicatablePost) (string, string, error) {
	if remoteID == "" {
		return "", "", fmt.Errorf("substack update: missing remote id")
	}
	if err := c.ensureProfile(ctx); err != nil {
		return "", "", err
	}
	var updated substackDraft
	if err := doJSON(ctx, c.http, http.MethodPut, c.baseURL+"/drafts/"+remoteID, c.headers(), c.draftPayload(post), &updated); err != nil {
		return "", "", fmt.Errorf("substack update draft: %w", err)
	}
	id := firstNonEmpty(updated.id(), remoteID)
	return id, firstNonEmpty(updated.publicURL(c.publicationURL), c.draftURL(id)), nil
}

func (c *SubstackClient) publishDraft(ctx context.Context, id string) (string, error) {
	body := map[string]any{"send": false, "share_automatically": false}
	var published substackDraft
	putURL := c.baseURL + "/drafts/" + id + "/publish"
	err := doJSON(ctx, c.http, http.MethodPut, putURL, c.headers(), body, &published)
	if err != nil && strings.Contains(err.Error(), "http 405") {
		err = doJSON(ctx, c.http, http.MethodPost, putURL, c.headers(), body, &published)
	}
	if err != nil {
		return "", err
	}
	return firstNonEmpty(published.publicURL(c.publicationURL), c.publicationURL+"/p/"+published.Slug), nil
}

func (c *SubstackClient) draftPayload(post model.SyndicatablePost) map[string]any {
	payload := map[string]any{
		"type":                      "newsletter",
		"draft_title":               post.Title,
		"draft_subtitle":            post.Description,
		"draft_body":                SubstackHTML(c.siteURL, CanonicalPost{Slug: post.Slug, Title: post.Title, ContentHTML: post.ContentHTML, CoverImage: post.CoverImage}),
		"draft_slug":                post.Slug,
		"audience":                  "everyone",
		"write_comment_permissions": "everyone",
		"search_engine_description": post.Description,
		"draft_bylines": []map[string]any{
			{"id": c.userID, "is_draft_byline": true},
		},
	}
	if c.pubUserID != 0 {
		payload["draft_bylines"] = []map[string]any{
			{"id": c.userID, "publicationUserId": c.pubUserID, "is_draft_byline": true},
		}
	}
	return payload
}

func (c *SubstackClient) ensureProfile(ctx context.Context) error {
	if c.userID != 0 {
		return nil
	}
	var profile substackProfile
	if err := doJSON(ctx, c.http, http.MethodGet, c.profileURL, c.headers(), nil, &profile); err != nil {
		return fmt.Errorf("substack profile: %w", err)
	}
	if profile.ID == 0 {
		return fmt.Errorf("substack profile: empty user id")
	}
	c.userID = profile.ID
	wantHost := hostLabel(c.publicationURL)
	for _, pu := range profile.PublicationUsers {
		host := pu.Publication.Subdomain + ".substack.com"
		if pu.Publication.CustomDomain != "" {
			host = pu.Publication.CustomDomain
		}
		if strings.EqualFold(host, wantHost) || strings.EqualFold(pu.Publication.Subdomain+".substack.com", wantHost) {
			c.pubUserID = pu.ID
			break
		}
	}
	if c.pubUserID == 0 && len(profile.PublicationUsers) > 0 {
		c.pubUserID = profile.PublicationUsers[0].ID
	}
	return nil
}

func (c *SubstackClient) headers() map[string]string {
	return map[string]string{"Cookie": c.cookie}
}

func (c *SubstackClient) draftURL(id string) string {
	return c.publicationURL + "/publish/post/" + id
}

type substackProfile struct {
	ID               int `json:"id"`
	PublicationUsers []struct {
		ID          int `json:"id"`
		Publication struct {
			Subdomain    string `json:"subdomain"`
			CustomDomain string `json:"custom_domain"`
		} `json:"publication"`
	} `json:"publicationUsers"`
}

type substackDraft struct {
	ID           json.RawMessage `json:"id"`
	Slug         string          `json:"slug"`
	CanonicalURL string          `json:"canonical_url"`
}

func (d substackDraft) id() string {
	raw := strings.TrimSpace(string(d.ID))
	return strings.Trim(raw, `"`)
}

func (d substackDraft) publicURL(publicationURL string) string {
	if d.CanonicalURL != "" {
		return d.CanonicalURL
	}
	if d.Slug != "" {
		return publicationURL + "/p/" + d.Slug
	}
	return ""
}

func substackCookieHeader(cookies, sid, connectSID string) string {
	if strings.TrimSpace(cookies) != "" {
		return strings.TrimSpace(cookies)
	}
	sid = cookieValue(sid)
	connectSID = cookieValue(connectSID)
	if sid == "" && connectSID == "" {
		return ""
	}
	if connectSID == "" {
		connectSID = sid
	}
	if sid == "" {
		sid = connectSID
	}
	return "substack.sid=" + sid + "; connect.sid=" + connectSID
}

func cookieValue(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if decoded, err := url.QueryUnescape(v); err == nil && decoded != "" {
		return decoded
	}
	return v
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
