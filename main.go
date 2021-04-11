package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	flagBase := flag.String("base", "dist", "the base directory")
	flagSerie := flag.String("serie", "-", "the serie")
	flagChapter := flag.String("chapter", "1", "the chapter")

	flag.Parse()

	serie := strings.ReplaceAll(*flagSerie, " ", "-")
	serie = strings.ReplaceAll(serie, "(", "")
	serie = strings.ReplaceAll(serie, ")", "")
	serie = strings.ToLower(serie)

	directory := *flagBase + "/" + *flagSerie + "/" + *flagChapter

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		log.Fatal(err)
	}

	collector := colly.NewCollector(
		// Only scrap the current page. Do not dive any deeper.
		colly.MaxDepth(1),
	)

	collector.OnRequest(func (request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.OnHTML("#all img", func (element *colly.HTMLElement) {
		fmt.Println("Downloading", element.Attr("data-src"))

		pageUrl := strings.Trim(element.Attr("data-src"), " ")
		segments := strings.Split(pageUrl, "/")
		filename := segments[len(segments) - 1]

		response, err := http.Get(pageUrl)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(directory + "/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(file, response.Body)
	})

	collector.Visit("https://readcomicsonline.ru/comic/" + serie + "/" + *flagChapter)
}
