package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/cache"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/handler"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/middleware"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/repository"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/service"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer rdb.Close()

	// Dependencies
	redisCache := cache.New(rdb)
	postRepo := repository.NewPostRepo(pool)
	tagRepo := repository.NewTagRepo(pool)
	markdownSvc := service.NewMarkdown()
	postSvc := service.NewPostService(postRepo, tagRepo, redisCache)
	syncSvc := service.NewSyncService(cfg, postRepo, tagRepo, markdownSvc, redisCache)

	// Handlers
	postHandler := handler.NewPostHandler(postSvc)
	tagHandler := handler.NewTagHandler(tagRepo, redisCache)
	syncHandler := handler.NewSyncHandler(syncSvc, cfg.SyncAPIKey)
	feedHandler := handler.NewFeedHandler(postSvc, cfg)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS(cfg.CORSOrigins))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handler.Health)
		r.Get("/posts", postHandler.List)
		r.Get("/posts/{slug}", postHandler.GetBySlug)
		r.Get("/tags", tagHandler.List)
		r.Get("/feed.xml", feedHandler.RSS)

		r.Group(func(r chi.Router) {
			r.Use(middleware.APIKey(cfg.SyncAPIKey))
			r.Post("/sync", syncHandler.Sync)
		})
	})

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
}
