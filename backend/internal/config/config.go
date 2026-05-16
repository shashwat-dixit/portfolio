package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	GitLabRepo    string
	GitLabToken   string
	SyncAPIKey    string
	CORSOrigins   []string
	SiteURL       string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/portfolio?sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		GitLabRepo:    getEnv("GITLAB_REPO", "https://gitlab.com/shashwat-dixit/blog.git"),
		GitLabToken:   getEnv("GITLAB_TOKEN", ""),
		SyncAPIKey:    getEnv("SYNC_API_KEY", ""),
		CORSOrigins:   []string{getEnv("CORS_ORIGIN", "http://localhost:4321")},
		SiteURL:       getEnv("SITE_URL", "https://shashwatdixit.com"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
