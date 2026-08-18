package model

import "time"

type Post struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ContentMD   string     `json:"contentMd,omitempty"`
	ContentHTML string     `json:"contentHtml,omitempty"`
	CoverImage  string     `json:"cover,omitempty"`
	Status      string     `json:"status"`
	ReadingTime int        `json:"readingTime"`
	Author      string     `json:"author"`
	PublishedAt *time.Time `json:"date"`
	UpdatedAt   *time.Time `json:"updated,omitempty"`
	GitLabSHA   string     `json:"-"`
	CreatedAt   time.Time  `json:"-"`
}

type PostSummary struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CoverImage  string     `json:"cover,omitempty"`
	Status      string     `json:"status"`
	ReadingTime int        `json:"readingTime"`
	Tags        []string   `json:"tags"`
	PublishedAt *time.Time `json:"date"`
}

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count int    `json:"count,omitempty"`
}

type PostListResponse struct {
	Posts      []PostSummary `json:"posts"`
	Pagination Pagination    `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type SyncResult struct {
	Synced   int              `json:"synced"`
	Created  int              `json:"created"`
	Updated  int              `json:"updated"`
	Deleted  int              `json:"deleted"`
	Medium   *CrossPostResult `json:"medium,omitempty"`
	Substack *CrossPostResult `json:"substack,omitempty"`
}

type CrossPostResult struct {
	Created int      `json:"created"`
	Updated int      `json:"updated"`
	Skipped int      `json:"skipped"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

const (
	PlatformMedium   = "medium"
	PlatformSubstack = "substack"
)

type Syndication struct {
	PostID      string
	Platform    string
	RemoteID    string
	RemoteURL   string
	ContentHash string
}

// SyndicatablePost is a published post plus the fields needed to cross-post it.
type SyndicatablePost struct {
	ID          string
	Slug        string
	Title       string
	Description string
	ContentMD   string
	ContentHTML string
	CoverImage  string
	Tags        []string
	ContentHash string
	PublishedAt *time.Time
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
