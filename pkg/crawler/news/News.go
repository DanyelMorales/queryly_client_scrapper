package news

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"golang.org/x/net/html"
	"news_scrapper/pkg/style"
	"strings"
	"time"
)

const Key = "ItemResult"

type Article struct {
	Id            string
	Author        string
	Content       string
	Title         string
	PublishedDate int64
	Keywords      []string
}

func CreateCollyCollector(context string) *colly.Collector {
	c := colly.NewCollector(
		colly.Async(true), // Enable asynchronous mode
	)

	err := c.Limit(&colly.LimitRule{
		Parallelism: 3,               // Max 3 requests at a time
		Delay:       1 * time.Second, // 1 second delay between requests,
		DomainGlob:  "*",
	})

	if err != nil {
		style.FailedAction(err.Error())
		return nil
	}

	c.OnError(func(response *colly.Response, err error) {
		style.FailedActionF("[%V] Error while scrapping website", context)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		r.Headers.Set("Referer", "https://google.com")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Cache-Control", "max-age=0")
	})
	return c
}

func StripHTML(htmlStr string) string {
	// Parse the HTML string
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return htmlStr
	}

	// Traverse the HTML node tree and extract text
	var output strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			output.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return output.String()
}
