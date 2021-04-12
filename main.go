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
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func download(base, serie, chapter string) {
	var files []string

	slug := strings.ReplaceAll(serie, " ", "-")
	slug = strings.ReplaceAll(slug, "(", "")
	slug = strings.ReplaceAll(slug, ")", "")
	slug = strings.ToLower(slug)

	// Pad the chapter with zeros so it's always 3-digits long. This helps when
	// storing .cbz files, so that readers correctly order them in the library.
	chapterPadded := fmt.Sprintf("%03v", chapter)

	// ---
	// Create the directory structure.

	directory := base + "/" + serie
	directoryTemp := "/tmp/" + serie + "/chapters/" + chapterPadded

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

	collector.Visit("https://readcomicsonline.ru/comic/" + slug + "/" + chapter)

	// ---
	// Create the .cbz file.

	filename := slug + "-" + chapterPadded + ".cbz"

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

	fmt.Println("Done", serie, chapter)
}

func main() {
	// ---
	// Define the CLI flags.

	base := flag.String("base", "dist", "the base directory")
	serie := flag.String("serie", "unknown", "the serie")
	chapter := flag.String("chapter", "1", "the chapter")

	flag.Parse()

	c := strings.Split(*chapter, "..")

	// ---
	// Download one chapter.

	if len(c) == 1 {
		download(*base, *serie, c[0])
	}

	// ---
	// Download all chapters.

	if len(c) == 2 {
		first, _ := strconv.Atoi(c[0])
		last, _ := strconv.Atoi(c[1])

		for i := first; i <= last; i++ {
			download(*base, *serie, strconv.Itoa(i))
		}
	}
}
