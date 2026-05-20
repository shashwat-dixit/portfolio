package model

import "time"

type Post struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ContentMD   string    `json:"-"`
	ContentHTML string    `json:"contentHtml,omitempty"`
	CoverImage  string    `json:"cover,omitempty"`
	Status      string    `json:"status"`
	ReadingTime int       `json:"readingTime"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"date"`
	UpdatedAt   time.Time `json:"updated,omitempty"`
	GitLabSHA   string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
}

type PostSummary struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover,omitempty"`
	Status      string    `json:"status"`
	ReadingTime int       `json:"readingTime"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"date"`
}

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count int    `json:"count,omitempty"`
}

type PostListResponse struct {
	Posts      []PostSummary `json:"posts"`
	Pagination Pagination   `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type SyncResult struct {
	Synced  int `json:"synced"`
	Created int `json:"created"`
	Updated int `json:"updated"`
	Deleted int `json:"deleted"`
}

// Frontmatter represents the YAML header in a blog markdown file.
type Frontmatter struct {
	Title       string   `yaml:"title"`
	Slug        string   `yaml:"slug"`
	Date        string   `yaml:"date"`
	Updated     string   `yaml:"updated,omitempty"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
	Cover       string   `yaml:"cover,omitempty"`
	Status      string   `yaml:"status"`
}
