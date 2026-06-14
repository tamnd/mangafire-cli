// Package mangafire scrapes the MangaFire.to manga catalog.
package mangafire

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Config controls the HTTP client behaviour.
type Config struct {
	BaseURL   string
	Rate      time.Duration
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://mangafire.to",
		Rate:      500 * time.Millisecond,
		Timeout:   30 * time.Second,
		Retries:   3,
		UserAgent: "Mozilla/5.0 (compatible; mangafire-cli/0.1; +https://github.com/tamnd/mangafire-cli)",
	}
}

// Client fetches MangaFire data.
type Client struct {
	cfg  Config
	http *http.Client
	last time.Time
}

// NewClient creates a Client from cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

var (
	itemRe = regexp.MustCompile(`(?s)href="(/manga/[^"]+)" class="poster"[^>]+>.*?<img[^>]+alt="([^"]+)"`)
	typeRe = regexp.MustCompile(`<span class="type">([^<]+)</span>`)
)

type filterResponse struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

// ListMangas fetches manga from the filter page.
// sort: most_viewed, latest_updated, new_release, title_az
// mangaType: "" (all), manga, manhwa, manhua, novel, one-shot, doujinshi
func (c *Client) ListMangas(ctx context.Context, sort, mangaType string, pages, limit int) ([]Manga, error) {
	if sort == "" {
		sort = "most_viewed"
	}
	var all []Manga
	rank := 1
	for page := 1; page <= pages || pages == 0; page++ {
		params := url.Values{
			"sort": {sort},
			"page": {fmt.Sprint(page)},
		}
		if mangaType != "" {
			params.Set("type[]", mangaType)
		}
		items, err := c.fetchFilterPage(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, it := range items {
			all = append(all, Manga{
				Rank:  rank,
				Slug:  strings.TrimPrefix(it[0], "/manga/"),
				Title: it[1],
				Type:  it[2],
				URL:   c.cfg.BaseURL + it[0],
			})
			rank++
			if limit > 0 && len(all) >= limit {
				return all, nil
			}
		}
		if len(items) == 0 || (limit > 0 && len(all) >= limit) {
			break
		}
	}
	return all, nil
}

func (c *Client) fetchFilterPage(ctx context.Context, params url.Values) ([][3]string, error) {
	path := "/filter?" + params.Encode()
	body, err := c.fetch(ctx, path, true)
	if err != nil {
		return nil, err
	}
	var resp filterResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing filter response: %w", err)
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("filter returned status %d", resp.Status)
	}
	html := resp.Result
	links := itemRe.FindAllStringSubmatch(html, -1)
	types := typeRe.FindAllStringSubmatch(html, -1)
	var out [][3]string
	for i, link := range links {
		t := ""
		if i < len(types) {
			t = types[i][1]
		}
		out = append(out, [3]string{link[1], link[2], t})
	}
	return out, nil
}

func (c *Client) fetch(ctx context.Context, path string, ajax bool) ([]byte, error) {
	reqURL := c.cfg.BaseURL + path
	var last error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		c.pace()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", c.cfg.UserAgent)
		if ajax {
			req.Header.Set("X-Requested-With", "XMLHttpRequest")
		}
		resp, err := c.http.Do(req)
		if err != nil {
			last = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			last = fmt.Errorf("HTTP %d", resp.StatusCode)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, reqURL)
		}
		return io.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("all retries failed for %s: %w", reqURL, last)
}

func (c *Client) pace() {
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}
