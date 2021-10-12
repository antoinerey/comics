package collector

import (
	"log"

	"github.com/gocolly/colly"
)

func CreateCollector() *colly.Collector {
	collector := colly.NewCollector(
		// Only scrap the current page. Do not dive any deeper.
		colly.MaxDepth(1),
	)

	collector.OnRequest(func(request *colly.Request) {
		log.Printf("Visiting %s", request.URL.String())
	})

	return collector
}
