package queryl

type ArticleHandler interface {
	Scrape()
	FetchMeta()
	Fetch()
	processArticles() func()
}
type Metadata struct {
	Query      string `json:"query"`
	Total      int    `json:"total"`
	EndIndex   int    `json:"endindex"`
	Correction string `json:"correction"`
	Suggest    string `json:"suggest"`
}

type FacetItem struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type Faceted struct {
	PubDate []FacetItem `json:"pubDate"`
	Creator []FacetItem `json:"creator"`
	Section []FacetItem `json:"section"`
}

type Item struct {
	ID          int    `json:"_id"`
	Index       int    `json:"index"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	PubDateUnix int64  `json:"pubdateunix"`
	Creator     string `json:"creator"`
	SubHeadline string `json:"subheadline"`
	Related     *[]string
}

type Response struct {
	Metadata Metadata `json:"metadata"`
	Faceted  Faceted  `json:"faceted"`
	Related  []string `json:"related"`
	Items    []Item   `json:"items"`
}

type SearchMode int

const (
	Query SearchMode = iota
	Section
)
