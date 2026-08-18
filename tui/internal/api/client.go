package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

type PostSummary struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ReadingTime int       `json:"readingTime"`
	Tags        []string  `json:"tags"`
	Date        *time.Time `json:"date"`
}

type Post struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ContentMD   string     `json:"contentMd"`
	Tags        []string   `json:"tags"`
	Date        *time.Time `json:"date"`
	ReadingTime int        `json:"readingTime"`
	Author      string     `json:"author"`
}

type listResponse struct {
	Posts      []PostSummary `json:"posts"`
	Pagination pagination    `json:"pagination"`
}

type pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 12 * time.Second},
	}
}

func (c *Client) ListPosts(ctx context.Context) ([]PostSummary, error) {
	var all []PostSummary
	page := 1
	for {
		u, err := url.Parse(c.BaseURL + "/api/posts")
		if err != nil {
			return nil, err
		}
		q := u.Query()
		q.Set("page", fmt.Sprintf("%d", page))
		q.Set("limit", "100")
		u.RawQuery = q.Encode()

		var resp listResponse
		if err := c.getJSON(ctx, u.String(), &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Posts...)
		if resp.Pagination.TotalPages == 0 || page >= resp.Pagination.TotalPages {
			break
		}
		page++
	}
	return all, nil
}

func (c *Client) GetPost(ctx context.Context, slug string) (*Post, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("slug required")
	}
	var post Post
	if err := c.getJSON(ctx, c.BaseURL+"/api/posts/"+url.PathEscape(slug), &post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (c *Client) getJSON(ctx context.Context, rawURL string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("decode api response: %w", err)
	}
	return nil
}

// StripFrontmatter removes YAML frontmatter so the body can be rendered as markdown.
func StripFrontmatter(md string) string {
	s := strings.TrimSpace(md)
	if !strings.HasPrefix(s, "---") {
		return md
	}
	rest := s[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return md
	}
	body := strings.TrimSpace(rest[idx+4:])
	if strings.HasPrefix(body, "---") {
		body = strings.TrimSpace(strings.TrimPrefix(body, "---"))
	}
	return body
}
