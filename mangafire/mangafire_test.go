package mangafire_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/mangafire-cli/mangafire"
)

const fakeHTML = `
<div class="unit item-1">
  <div class="inner">
    <a href="/manga/one-piecee.dkw" class="poster" data-tip="1?/x">
      <div><img src="img.jpg" alt="One Piece"></div>
    </a>
    <div class="info">
      <div><span class="type">Manga</span></div>
      <a href="/manga/one-piecee.dkw">One Piece</a>
    </div>
  </div>
</div>
<div class="unit item-2">
  <div class="inner">
    <a href="/manga/naruto.abc" class="poster" data-tip="2?/x">
      <div><img src="img2.jpg" alt="Naruto"></div>
    </a>
    <div class="info">
      <div><span class="type">Manga</span></div>
      <a href="/manga/naruto.abc">Naruto</a>
    </div>
  </div>
</div>
`

func newTestClient(ts *httptest.Server) *mangafire.Client {
	cfg := mangafire.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return mangafire.NewClient(cfg)
}

func TestListMangas(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
			t.Error("expected XMLHttpRequest header")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200,
			"result": fakeHTML,
		})
	}))
	defer ts.Close()

	c := newTestClient(ts)
	items, err := c.ListMangas(context.Background(), "most_viewed", "", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 manga, got %d", len(items))
	}
	if items[0].Title != "One Piece" {
		t.Errorf("first title = %q, want One Piece", items[0].Title)
	}
	if items[0].Slug != "one-piecee.dkw" {
		t.Errorf("first slug = %q, want one-piecee.dkw", items[0].Slug)
	}
	if items[0].Type != "Manga" {
		t.Errorf("first type = %q, want Manga", items[0].Type)
	}
}
