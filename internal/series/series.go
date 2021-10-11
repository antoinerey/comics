package series

import (
	"github.com/gocolly/colly"
)

var collector = colly.NewCollector(
	// Only scrap the current page. Do not dive any deeper.
	colly.MaxDepth(1),
)

type Series = struct {
	Url        string
	Name       string
	IssuesURLs []string
}

func ParseURL(url string) (Series, error) {
	var name string
	var issuesURLs []string

	collector.OnHTML(".anime-details .title", func(element *colly.HTMLElement) {
		name = element.Text
	})

	collector.OnHTML(".basic-list a", func(element *colly.HTMLElement) {
		// Prepend the URL to the list.
		issuesURLs = append([]string{element.Attr("href")}, issuesURLs...)
	})

	err := collector.Visit(url)
	if err != nil {
		return Series{}, err
	}

	series := Series{
		Url:        url,
		Name:       name,
		IssuesURLs: issuesURLs,
	}

	return series, nil
}
