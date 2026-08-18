package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
)

const mediumAPI = "https://api.medium.com/v1"

type MediumClient struct {
	token         string
	publicationID string
	publishStatus string
	siteURL       string
	baseURL       string
	http          httpDoer
	userID        string
}

func NewMediumClient(cfg *config.Config) *MediumClient {
	return &MediumClient{
		token:         cfg.MediumToken,
		publicationID: cfg.MediumPublicationID,
		publishStatus: cfg.MediumPublishStatus,
		siteURL:       cfg.SiteURL,
		baseURL:       mediumAPI,
		http:          defaultHTTPClient(),
	}
}

func (c *MediumClient) Enabled() bool {
	return c != nil && strings.TrimSpace(c.token) != ""
}

func (c *MediumClient) SupportsUpdate() bool { return false }

func (c *MediumClient) Publish(ctx context.Context, post model.SyndicatablePost) (remoteID, remoteURL string, err error) {
	userID, err := c.me(ctx)
	if err != nil {
		return "", "", err
	}

	payload := map[string]any{
		"title":           truncateRunes(post.Title, 100),
		"contentFormat":   "markdown",
		"content":         MediumContent(c.siteURL, CanonicalPost{Slug: post.Slug, Title: post.Title, ContentMD: post.ContentMD, CoverImage: post.CoverImage}),
		"tags":            MediumTags(post.Tags),
		"canonicalUrl":    CanonicalPostURL(c.siteURL, post.Slug),
		"publishStatus":   c.publishStatus,
		"notifyFollowers": false,
	}

	endpoint := fmt.Sprintf("%s/users/%s/posts", c.baseURL, userID)
	if c.publicationID != "" {
		endpoint = fmt.Sprintf("%s/publications/%s/posts", c.baseURL, c.publicationID)
	}

	var resp struct {
		Data struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := doJSON(ctx, c.http, http.MethodPost, endpoint, c.headers(), payload, &resp); err != nil {
		return "", "", fmt.Errorf("medium create post: %w", err)
	}
	if resp.Data.ID == "" {
		return "", "", fmt.Errorf("medium create post: empty id")
	}
	return resp.Data.ID, resp.Data.URL, nil
}

func (c *MediumClient) Update(_ context.Context, _ string, _ model.SyndicatablePost) (string, string, error) {
	return "", "", fmt.Errorf("medium API cannot update existing posts")
}

func (c *MediumClient) me(ctx context.Context) (string, error) {
	if c.userID != "" {
		return c.userID, nil
	}
	var resp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := doJSON(ctx, c.http, http.MethodGet, c.baseURL+"/me", c.headers(), nil, &resp); err != nil {
		return "", fmt.Errorf("medium me: %w", err)
	}
	if resp.Data.ID == "" {
		return "", fmt.Errorf("medium me: empty user id")
	}
	c.userID = resp.Data.ID
	return c.userID, nil
}

func (c *MediumClient) headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + c.token,
	}
}
