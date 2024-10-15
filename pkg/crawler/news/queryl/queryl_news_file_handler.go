package queryl

import (
	"fmt"
	"news_scrapper/pkg/crawler/news"
	"news_scrapper/pkg/style"
)

type Handler struct {
	Client *Queryl
}

var facetedFile string

type BatchAction func(handler Handler)

type FakingTheFunk struct {
	Enabled                         bool
	Context, Host, ApiKey, Selector string
	OverrideExistingNews            bool
}

var ValidSortDate = map[string]int{
	"last-week":  168,
	"today":      24,
	"last-month": 720,
}

func TriggerBatchScrapping(opts []FakingTheFunk, action BatchAction) {
	for i := range opts {
		instance := opts[i]
		if instance.Enabled {
			handler := NewQuerylFileHandler(instance.Context, instance.Host, instance.ApiKey, instance.Selector, instance.OverrideExistingNews)
			handler.Client.FullLog = true
			handler.Client.SaveData = true
			action(handler)
		}
	}
}

func NewQuerylFileHandler(context, host, apiKey, selector string, overrideExistingNews bool) Handler {
	articleChan := make(chan news.Article)
	articleStreamEnd := make(chan bool)
	client := NewQuerylClient(context, host, apiKey, selector, articleChan, articleStreamEnd)
	client.Defaults(overrideExistingNews, "json")
	facetedFile = fmt.Sprintf(client.OutputFile, "faceted")
	return Handler{Client: client}
}

func (q Handler) Scrape(mode SearchMode, value string, dateSort string, batchSize int) {
	var dateSortVal = 0
	if value, exists := ValidSortDate[dateSort]; exists {
		dateSortVal = value
	} else {
		style.ExitActionF("invalid sort value: %v", dateSort)
		return
	}
	go q.processArticles()()
	q.Client.Scrape(mode, value, dateSortVal, batchSize)
}

func (q Handler) FetchMeta() (*Faceted, error) {
	return q.Client.LoadFaceted(facetedFile)
}

func (q Handler) processArticles() func() {
	return func() {
		news.ArticleWatcher(q.Client.Client)
	}
}
