package service

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
)

type fakePosts struct {
	posts []model.SyndicatablePost
	err   error
}

func (f *fakePosts) ListPublishedFull(context.Context, int) ([]model.SyndicatablePost, error) {
	return f.posts, f.err
}

type fakeSyn struct {
	existing map[string]map[string]model.Syndication
	upserts  []model.Syndication
}

func (f *fakeSyn) List(context.Context) (map[string]map[string]model.Syndication, error) {
	if f.existing == nil {
		return map[string]map[string]model.Syndication{}, nil
	}
	return f.existing, nil
}

func (f *fakeSyn) Upsert(_ context.Context, s model.Syndication) error {
	f.upserts = append(f.upserts, s)
	return nil
}

type fakePoster struct {
	updateOK  bool
	publishes int
	updates   int
	failSlug  string
}

func (f *fakePoster) Enabled() bool        { return true }
func (f *fakePoster) SupportsUpdate() bool { return f.updateOK }
func (f *fakePoster) Publish(_ context.Context, post model.SyndicatablePost) (string, string, error) {
	f.publishes++
	if post.Slug == f.failSlug {
		return "", "", fmt.Errorf("boom")
	}
	return "id-" + post.Slug, "https://example.com/" + post.Slug, nil
}
func (f *fakePoster) Update(_ context.Context, remoteID string, post model.SyndicatablePost) (string, string, error) {
	f.updates++
	return remoteID, "https://example.com/" + post.Slug, nil
}

func TestSyndicateCreatesSkipsAndUpdates(t *testing.T) {
	posts := []model.SyndicatablePost{
		{ID: "1", Slug: "new", Title: "New", ContentHash: "aaa"},
		{ID: "2", Slug: "same", Title: "Same", ContentHash: "bbb"},
		{ID: "3", Slug: "changed", Title: "Changed", ContentHash: "ccc2"},
	}
	existing := map[string]map[string]model.Syndication{
		"2": {model.PlatformMedium: {PostID: "2", Platform: model.PlatformMedium, RemoteID: "m2", ContentHash: "bbb"}},
		"3": {model.PlatformMedium: {PostID: "3", Platform: model.PlatformMedium, RemoteID: "m3", ContentHash: "ccc1"}},
	}

	svc := &SyndicationService{
		cfg:      &config.Config{},
		postRepo: &fakePosts{posts: posts},
		synRepo:  &fakeSyn{existing: existing},
	}

	noUpdate := &fakePoster{}
	created := svc.syndicatePlatform(context.Background(), model.PlatformMedium, noUpdate, posts, existing)
	if created.Created != 1 || created.Skipped != 2 || created.Failed != 0 {
		t.Fatalf("no-update result: %+v publishes=%d", created, noUpdate.publishes)
	}

	store := &fakeSyn{existing: existing}
	svc.synRepo = store
	updater := &fakePoster{updateOK: true}
	updated := svc.syndicatePlatform(context.Background(), model.PlatformMedium, updater, posts, existing)
	if updated.Created != 1 || updated.Updated != 1 || updated.Skipped != 1 {
		t.Fatalf("update result: %+v publishes=%d updates=%d", updated, updater.publishes, updater.updates)
	}
	if len(store.upserts) != 2 {
		t.Fatalf("upserts %d", len(store.upserts))
	}
}

func TestSyndicateRecordsFailures(t *testing.T) {
	posts := []model.SyndicatablePost{{ID: "1", Slug: "bad", ContentHash: "x"}}
	svc := &SyndicationService{postRepo: &fakePosts{posts: posts}, synRepo: &fakeSyn{}}
	poster := &fakePoster{failSlug: "bad"}
	got := svc.syndicatePlatform(context.Background(), model.PlatformMedium, poster, posts, map[string]map[string]model.Syndication{})
	if got.Failed != 1 || got.Created != 0 || len(got.Errors) != 1 {
		t.Fatalf("%+v", got)
	}
}

func TestSyndicationDisabledWithoutCredentials(t *testing.T) {
	svc := NewSyndicationService(&config.Config{}, nil, nil)
	if svc.Enabled() {
		t.Fatal("expected disabled")
	}
	m, s := svc.Syndicate(context.Background())
	if m != nil || s != nil {
		t.Fatalf("got %+v %+v", m, s)
	}
}
