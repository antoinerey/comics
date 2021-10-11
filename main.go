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

func download(directory, series, issue string) {
	var imgs []*colly.HTMLElement
	var files []string

	slug := strings.ReplaceAll(series, " ", "-")
	slug = strings.ReplaceAll(slug, ":", "")
	slug = strings.ReplaceAll(slug, "(", "")
	slug = strings.ReplaceAll(slug, ")", "")
	slug = strings.ToLower(slug)

	// Pad the issue with zeros so it's always 3-digits long. This helps when
	// storing .cbz files, so that readers correctly order them in the library.
	issuePadded := fmt.Sprintf("%03v", issue)

	// ---
	// Download each images.

	collector := colly.NewCollector(
		// Only scrap the current page. Do not dive any deeper.
		colly.MaxDepth(1),
	)

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.OnHTML(".chapter-container > img", func(element *colly.HTMLElement) {
		imgs = append(imgs, element)
	})

	collector.Visit("https://viewcomics.me/" + slug + "/issue-" + issue + "/full")

	if len(imgs) == 0 {
		log.Fatal("No images found.")
	}

	// ---
	// Create the directory structure.

	directorySeries := directory + "/" + series
	directorySeriesTemp := "/tmp/" + series
	directoryIssueTemp := directorySeriesTemp + "/issues/" + issuePadded
	defer os.RemoveAll(directorySeriesTemp)

	err := os.MkdirAll(directorySeries, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(directoryIssueTemp, 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range imgs {
		fmt.Println("Downloading", img.Attr("alt"))

		pageUrl := strings.Trim(img.Attr("src"), " ")
		filename := img.Attr("alt") + ".jpg"

		response, err := http.Get(pageUrl)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(directoryIssueTemp + "/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(file, response.Body)
		files = append(files, filename)
	}

	// ---
	// Create the .cbz file.

	filename := slug + "-" + issuePadded + ".cbz"

	zipFile, err := os.Create(directorySeries + "/" + filename)
	if err != nil {
		log.Fatal(err)
	}

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		content, err := ioutil.ReadFile(directoryIssueTemp + "/" + file)
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

	fmt.Println("Done", series, issue)
}

func main() {
	// ---
	// Define the CLI flags.

	directory := flag.String("directory", "dist", "")
	series := flag.String("series", "unknown", "")
	issues := flag.String("issues", "1", "")

	flag.Parse()

	// i := strings.Split(*issues, "..")

	// ---
	// Download one issue.

	// if len(i) == 1 {
	// 	download(*directory, *series, i[0])
	// }

	// ---
	// Download all issues.

	// if len(i) == 2 {
	// 	first, err := strconv.Atoi(i[0])
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	last, err := strconv.Atoi(i[1])
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	for i := first; i <= last; i++ {
	// 		download(*directory, *series, strconv.Itoa(i))
	// 	}
	// }

	log.Print(*series)

	collector := colly.NewCollector(
		// Only scrap the current page. Do not dive any deeper.
		colly.MaxDepth(1),
	)

	collector.OnHTML(".basic-list a", func(element *colly.HTMLElement) {
		log.Print(element.Attr("href"))
	})

	collector.Visit(*series)
}
