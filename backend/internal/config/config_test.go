package config

import (
	"testing"
)

func TestNormalizePublishStatus(t *testing.T) {
	tests := map[string]string{
		"draft":    "draft",
		"UNLISTED": "unlisted",
		"":         "public",
		"public":   "public",
		"other":    "public",
	}
	for in, want := range tests {
		if got := normalizePublishStatus(in); got != want {
			t.Fatalf("%q: got %q want %q", in, got, want)
		}
	}
}

func TestGetEnvBool(t *testing.T) {
	t.Setenv("X_BOOL", "true")
	if !getEnvBool("X_BOOL", false) {
		t.Fatal("true")
	}
	t.Setenv("X_BOOL", "nope")
	if getEnvBool("X_BOOL", true) {
		t.Fatal("invalid should be false")
	}
	t.Setenv("X_BOOL", "")
	if !getEnvBool("X_BOOL", true) {
		t.Fatal("empty uses fallback")
	}
}
