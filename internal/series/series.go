package series

import (
	"fmt"
	"strings"

	"github.com/antoinerey/comics/internal/collector"
	"github.com/antoinerey/comics/internal/issues"
	"github.com/gocolly/colly"
)

type Series struct {
	URL    string
	Title  string
	Issues []issues.Issue

	collector *colly.Collector
}

func CreateSeries(URL string) Series {
	return Series{
		URL:       URL,
		collector: collector.CreateCollector(),
	}
}

func (series Series) Parse() (Series, error) {
	series.collector.OnHTML(".anime-details .title", func(element *colly.HTMLElement) {
		// The parsed HTML includes new line character.
		series.Title = strings.Trim(element.Text, "\n")
	})

	series.collector.OnHTML(".basic-list a", func(element *colly.HTMLElement) {
		// Save the full issue URL.
		issueURL := fmt.Sprintf("%s/full", element.Attr("href"))
		// Prepend the URL to the list.
		series.Issues = append([]issues.Issue{issues.CreateIssue(issueURL)}, series.Issues...)
	})

	err := series.collector.Visit(series.URL)

	return series, err
}

func (series Series) GetDirectory(root string) string {
	return fmt.Sprintf("%s/%s", root, series.Title)
}
