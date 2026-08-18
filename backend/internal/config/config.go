package config

import (
	"os"
	"strings"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RedisTLS      bool
	GitLabRepo    string
	GitLabToken   string
	SyncAPIKey    string
	CORSOrigins   []string
	SiteURL       string

	MediumToken         string
	MediumPublicationID string
	MediumPublishStatus string

	SubstackPublicationURL string
	SubstackSID            string
	SubstackConnectSID     string
	SubstackCookies        string
	SubstackPublish        bool
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/portfolio?sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		RedisTLS:      getEnv("REDIS_TLS", "") == "true",
		GitLabRepo:    getEnv("GITLAB_REPO", "https://gitlab.com/shashwat-dixit/blog.git"),
		GitLabToken:   getEnv("GITLAB_TOKEN", ""),
		SyncAPIKey:    getEnv("SYNC_API_KEY", ""),
		CORSOrigins:   []string{getEnv("CORS_ORIGIN", "http://localhost:4321")},
		SiteURL:       getEnv("SITE_URL", "https://shashwatdixit.com"),

		MediumToken:         getEnv("MEDIUM_TOKEN", ""),
		MediumPublicationID: getEnv("MEDIUM_PUBLICATION_ID", ""),
		MediumPublishStatus: normalizePublishStatus(getEnv("MEDIUM_PUBLISH_STATUS", "public")),

		SubstackPublicationURL: strings.TrimRight(getEnv("SUBSTACK_PUBLICATION_URL", ""), "/"),
		SubstackSID:            getEnv("SUBSTACK_SID", ""),
		SubstackConnectSID:     getEnv("SUBSTACK_CONNECT_SID", ""),
		SubstackCookies:        getEnv("SUBSTACK_COOKIES", ""),
		SubstackPublish:        getEnvBool("SUBSTACK_PUBLISH", false),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return fallback
	}
	return v == "true" || v == "1" || v == "yes"
}

func normalizePublishStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "draft":
		return "draft"
	case "unlisted":
		return "unlisted"
	default:
		return "public"
	}
}
