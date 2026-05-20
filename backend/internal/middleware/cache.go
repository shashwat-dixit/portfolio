package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"strings"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/cache"
)

func CacheAside(rc *cache.RedisCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			key := cacheKeyFromRequest(r)
			ctx := r.Context()

			cached, err := rc.Get(ctx, key)
			if err != nil {
				slog.Warn("cache middleware get error", "key", key, "error", err)
			}
			if cached != "" {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Cache", "HIT")
				w.Write([]byte(cached))
				return
			}

			rec := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}, statusCode: http.StatusOK}
			next.ServeHTTP(rec, r)

			if rec.statusCode >= 200 && rec.statusCode < 300 {
				if err := rc.Set(ctx, key, rec.body.String()); err != nil {
					slog.Warn("cache middleware set error", "key", key, "error", err)
				}
			}
			w.Header().Set("X-Cache", "MISS")
		})
	}
}

func cacheKeyFromRequest(r *http.Request) string {
	key := cache.KeyPrefix + "http:" + r.URL.Path
	if q := r.URL.RawQuery; q != "" {
		key += "?" + q
	}
	return strings.ReplaceAll(key, " ", "")
}

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
