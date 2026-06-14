package mangafire

// Manga is a title entry from MangaFire.
type Manga struct {
	Rank  int    `json:"rank"  csv:"rank"  tsv:"rank"`
	Slug  string `json:"slug"  csv:"slug"  tsv:"slug"`
	Title string `json:"title" csv:"title" tsv:"title"`
	Type  string `json:"type"  csv:"type"  tsv:"type"`
	URL   string `json:"url"   csv:"url"   tsv:"url"`
}
