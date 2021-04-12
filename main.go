package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	var files []string

	// ---
	// Define the CLI flags.

	flagBase := flag.String("base", "dist", "the base directory")
	flagSerie := flag.String("serie", "-", "the serie")
	flagChapter := flag.String("chapter", "1", "the chapter")

	flag.Parse()

	// ---
	// Create the directory structure.

	serie := strings.ReplaceAll(*flagSerie, " ", "-")
	serie = strings.ReplaceAll(serie, "(", "")
	serie = strings.ReplaceAll(serie, ")", "")
	serie = strings.ToLower(serie)

	directory := *flagBase + "/" + *flagSerie
	directoryTemp := "/tmp/" + *flagSerie + "/chapters/" + *flagChapter

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(directoryTemp, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// ---
	// Download each images.

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

		file, err := os.Create(directoryTemp + "/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(file, response.Body)
		files = append(files, filename)
	})

	collector.Visit("https://readcomicsonline.ru/comic/" + serie + "/" + *flagChapter)

	// ---
	// Create the .cbz file.

	filename := serie + "-" + *flagChapter + ".cbz"

	zipFile, err := os.Create(directory + "/" + filename)
	if err != nil {
		log.Fatal(err)
	}

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		content, err := ioutil.ReadFile(directoryTemp + "/" + file)
		if err != nil {
			log.Fatal(err)
		}

		fileWriter, err := zipWriter.Create(file)
		if err != nil {
			log.Fatal(err)
		}

		fileWriter.Write([]byte(string(content)))
	}

	// ---
	// All good.

	fmt.Println("Done")
}
