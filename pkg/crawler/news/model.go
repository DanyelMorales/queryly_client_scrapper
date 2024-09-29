package news

import (
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
)

const OutputPath = "outputPath"
const SubDirDefined = "out"

type ShouldScrapeLink func(article string) bool

type Client struct {
	ArticleChan          chan Article
	ArticleStreamEnd     chan bool
	FullLog              bool
	SaveData             bool
	Context              string
	OutputFile           string
	DirName              string
	OutputFileExt        string
	ApiKey               string
	Resty                *resty.Client
	Colly                *colly.Collector
	FullHost             string
	ShouldScrapeLink     ShouldScrapeLink
	InMemoryExistingNews []string
}

func (q *Client) Defaults(overrideExistingNews bool, outputExtension string) {
	q.OutputFileExt = outputExtension
	dirName, outputFile := GetDirNameAndOutputName(q.Context, outputExtension)
	q.DirName = dirName
	q.OutputFile = outputFile
	q.FullLog = true
	q.SaveData = true
	if overrideExistingNews {
		q.InMemoryExistingNews = FetchExistingNews(q.DirName)
		q.ShouldScrapeLink = q.ShouldPerformOperationOnArticle
	}
}
func (q Client) ShouldPerformOperationOnArticle(article string) bool {
	return !ArticleExists(&q.InMemoryExistingNews, article+q.OutputFileExt)
}
