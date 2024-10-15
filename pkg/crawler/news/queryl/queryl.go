package queryl

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
	"news_scrapper/pkg/crawler/news"
	"news_scrapper/pkg/style"
	"strconv"
)

const SearchHost = "https://api.queryly.com"
const SearchUrl = SearchHost + "/json.aspx"

type Queryl struct {
	*news.Client
}

func NewQuerylClient(context, fullHost, apiKey, selector string, articleChan chan news.Article, articleStreamEnd chan bool) *Queryl {
	c := news.CreateCollyCollector(context)
	if c == nil {
		return nil
	}
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		var articleContent string
		e.DOM.Find("script").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		e.DOM.Find("meta").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		e.DOM.Each(func(i int, s *goquery.Selection) {
			articleContent += s.Text()
		})
		articleData := e.Request.Ctx.GetAny(news.Key).(Item)
		article := news.Article{
			Id:            strconv.Itoa(articleData.ID),
			Author:        articleData.Creator,
			Content:       articleContent,
			Title:         articleData.Title,
			PublishedDate: articleData.PubDateUnix,
			Keywords:      *articleData.Related,
		}
		articleChan <- article
		style.SuccessfulActionF("[%V] article sent to channel, ArticleID:%V", context, articleData.ID)
	})
	cl := &news.Client{FullHost: fullHost, Colly: c, ApiKey: apiKey, Resty: resty.New(), ArticleChan: articleChan, Context: context, ArticleStreamEnd: articleStreamEnd}
	return &Queryl{Client: cl}
}

func (newsClient *Queryl) Scrape(mode SearchMode, value string, dateSort int, batchSize int) {
	value = ProcessString(value)
	result, err := newsClient.Fetch(mode, value, dateSort, batchSize)
	if err != nil {
		style.FailedActionF("[%V][Scrape] error detected %v", newsClient.Context, err.Error())
		return
	}
	newsClient.extractArticleContent(result)
}

func (newsClient *Queryl) Fetch(mode SearchMode, value string, dateSort int, batchSize int) ([]Item, error) {
	switch mode {
	case Query:
		return newsClient.fetchArticlesByTerm(value, dateSort, batchSize)
	case Section:
		return newsClient.fetchArticlesBySection(value, dateSort, batchSize)
	default:
		return newsClient.fetchArticlesBySection(value, dateSort, batchSize)
	}
}

func (newsClient *Queryl) fetchMeta() (*Faceted, error) {
	queryParams := BuildQueryParams(newsClient.ApiKey, strconv.Itoa(0), "1")
	(*queryParams)["query"] = "1"
	return newsClient.fetchFaceted(*queryParams)
}

func (newsClient *Queryl) fetchArticlesByTerm(term string, dateSort int, batchSize int) ([]Item, error) {
	endindex := viper.GetString("endIndex")
	if endindex == "" {
		endindex = "0"
	}
	queryParams := BuildQueryParams(newsClient.ApiKey, endindex, strconv.Itoa(batchSize))
	newsClient.appendFacet("pubDate", strconv.Itoa(dateSort), queryParams)
	newsClient.appendFacet("query", "", queryParams)
	(*queryParams)["query"] = term
	(*queryParams)["sort"] = "date"
	return newsClient.fetchArticleLinks(*queryParams)
}

func (newsClient *Queryl) extractArticleContent(results []Item) {
	for i := range results {
		link := newsClient.FullHost + results[i].Link
		ctx := colly.NewContext()
		ctx.Put(news.Key, results[i])
		if newsClient.ShouldScrapeLink != nil && !newsClient.ShouldScrapeLink(strconv.Itoa(results[i].ID)) {
			style.OnNoticeF("[%V] this article already exists, skipping scrapping  %v", newsClient.Context, link)
			continue
		} else {
			style.OnNoticeF("[%V] this article doesn't exist in local, proceeding to scrapping %v", newsClient.Context, link)
		}
		err := newsClient.Colly.Request("GET", link, nil, ctx, nil)
		if err != nil {
			style.FailedActionF("[%V] failed requesting the following link %v", newsClient.Context, link)
		} else {
			style.OnActionF("[%V] about to scrape %V", newsClient.Context, link)
		}
	}
	newsClient.Colly.Wait()
	newsClient.ArticleStreamEnd <- true
}

func (newsClient *Queryl) fetchArticlesBySection(section string, dateSort int, batchSize int) ([]Item, error) {
	queryParams := BuildQueryParams(newsClient.ApiKey, strconv.Itoa(0), strconv.Itoa(batchSize))
	newsClient.appendFacet("section", section, queryParams)
	newsClient.appendFacet("pubDate", strconv.Itoa(dateSort), queryParams)
	(*queryParams)["sort"] = "date"
	style.OnActionF("[%V] about to invoke Article endpoint", newsClient.Context)
	return newsClient.fetchArticleLinks(*queryParams)
}

func (newsClient *Queryl) appendFacet(facetedKey string, facetedValue string, queryParams *map[string]string) {
	if val, exists := (*queryParams)["facetedkey"]; exists {
		val += facetedKey
		(*queryParams)["facetedkey"] = val + "|"
	} else {
		(*queryParams)["facetedkey"] = facetedKey + "|"
	}
	if val, exists := (*queryParams)["facetedvalue"]; exists {
		val += facetedValue
		(*queryParams)["facetedvalue"] = val + "|"
	} else {
		(*queryParams)["facetedvalue"] = facetedValue + "|"
	}
}

func (newsClient *Queryl) fetchFaceted(queryParams map[string]string) (*Faceted, error) {
	resp, err := newsClient.Resty.R().SetQueryParams(queryParams).SetHeader("Accept", "application/json").Get(SearchUrl)
	if err != nil {
		style.FailedActionF("[%V][fetchFaceted] msg=error occurred while invoking %V, error=%v", newsClient.Context, SearchUrl, err.Error())
		return nil, err
	}
	var result Faceted
	var responseBody Response
	if resp.StatusCode() == 200 {
		err = json.Unmarshal(resp.Body(), &responseBody)
		if err != nil {
			style.FailedActionF("[%V][fetchFaceted] unmarshalling failed with %v", newsClient.Context, err.Error())
			return nil, err
		} else {
			result = responseBody.Faceted
		}
	} else {
		style.FailedActionF("[%V][fetchFaceted] bad status code  %v", newsClient.Context, resp.StatusCode())
		return nil, err
	}
	style.SuccessfulActionF("[%V][fetchFaceted] successfully received response for %V, items:%V", newsClient.Context, SearchUrl, len(responseBody.Faceted.Section))
	return &result, nil
}

func (newsClient *Queryl) fetchArticleLinks(queryParams map[string]string) ([]Item, error) {
	resp, err := newsClient.Resty.R().SetQueryParams(queryParams).SetHeader("Accept", "application/json").Get(SearchUrl)
	if err != nil {
		style.FailedActionF("[%V][fetchArticleLinks] msg=error occurred while invoking %V, error=%v", newsClient.Context, SearchUrl, err.Error())
		return nil, err
	}
	var result []Item
	var responseBody Response
	if resp.StatusCode() == 200 {
		err = json.Unmarshal(resp.Body(), &responseBody)
		if err != nil {
			style.FailedActionF("[%V][fetchArticleLinks] unmarshalling failed with %v", newsClient.Context, err.Error())
			return nil, err
		} else {
			result = responseBody.Items
			for i := range result {
				result[i].Related = &responseBody.Related
			}
		}
	} else {
		style.FailedActionF("[%V][fetchArticleLinks] bad status code  %v", newsClient.Context, resp.StatusCode())
		return nil, err
	}
	style.SuccessfulActionF("[%V][fetchArticleLinks] successfully received response for %V, items=%V", newsClient.Context, SearchUrl, len(result))
	return result, nil
}

func (q Queryl) LoadFaceted(facetedFile string) (*Faceted, error) {
	if !news.FileExists(facetedFile) {
		style.OnActionF("[%V] Faceted is not in local, attempting to fetch it from server: %V", q.Context, facetedFile)
		meta, err := q.fetchMeta()
		if err != nil {
			return nil, err
		}
		jsonData, err := json.Marshal(meta)
		if err != nil {
			style.FailedActionF("Error converting struct to JSON:%V", err)
			return nil, err
		}
		style.OkLogCF(true, "[%V] Received Faceted, about to write the content into file: %V", q.Context, facetedFile)
		err = news.WriteFile(facetedFile, string(jsonData))
		if err != nil {
			style.FailedActionF("Error writing file:%V", err)
			return nil, err
		}
		style.SuccessfulActionF("[%V] Done, returning facet now: %V", q.Context, facetedFile)
		return meta, nil
	}
	file, err := news.ReadFile(facetedFile)
	if err != nil {
		return nil, err
	}
	var faceted Faceted
	err = json.Unmarshal([]byte(file), &faceted)
	if err != nil {
		return nil, err
	}
	style.OnNoticeF("[%V] found Facet locally returning facet now: %V", q.Context, facetedFile)
	return &faceted, nil
}
