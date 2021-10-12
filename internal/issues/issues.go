package issues

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antoinerey/comics/internal/collector"
	"github.com/gocolly/colly"
)

type Issue struct {
	URL   string
	Title string
	Pages []*colly.HTMLElement

	collector *colly.Collector
}

func CreateIssue(URL string) Issue {
	return Issue{
		URL:       URL,
		collector: collector.CreateCollector(),
	}
}

func (issue Issue) Parse() Issue {
	issue.collector.OnHTML(".chapter-title h1", func(element *colly.HTMLElement) {
		// The parsed HTML includes new line character.
		issue.Title = strings.Trim(element.Text, "\n")
	})

	issue.collector.OnHTML(".chapter-container > img", func(page *colly.HTMLElement) {
		issue.Pages = append(issue.Pages, page)
	})

	err := issue.collector.Visit(issue.URL)
	if err != nil {
		log.Printf("Failed to visit %s", issue.URL)
		log.Fatal(err)
	}

	return issue
}

func (issue Issue) Download(baseDir, tmpDir string) {
	err := os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)
	if err != nil {
		log.Printf("Failed to create temporary directory at path %s", tmpDir)
		log.Fatal(err)
	}

	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		log.Printf("Failed to create library directory at path %s", baseDir)
		log.Fatal(err)
	}

	zipPath := fmt.Sprintf("%s/%s.cbz", tmpDir, issue.Title)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Printf("Failed to create .cbz file at path %s", zipPath)
		log.Fatal(err)
	}

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, page := range issue.Pages {
		src := strings.Trim(page.Attr("src"), " ")
		title := page.Attr("alt") + ".jpg"

		log.Printf("Downloading page %s", title)

		response, err := http.Get(src)
		if err != nil {
			log.Printf("Failed to request the given page URL %s", src)
			log.Fatal(err)
		}

		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print("Failed to read response body")
			log.Fatal(err)
		}

		fileWriter, err := zipWriter.Create(title)
		if err != nil {
			log.Printf("Failed to create zip writer for %s", title)
			log.Fatal(err)
		}

		_, err = fileWriter.Write(content)
		if err != nil {
			log.Print("Failed to write .cbz file")
			log.Fatal(err)
		}
	}

	err = os.Rename(zipPath, fmt.Sprintf("%s/%s.cbz", baseDir, issue.Title))
	if err != nil {
		log.Print("Failed to move temporary .cbz file to library path")
		log.Fatal(err)
	}

	log.Printf("Successfuly downloaded issue %s", issue.Title)
}

func (issue Issue) IsMissing(path string) bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s.cbz", path, issue.Title))

	if err == nil {
		return false
	}

	if os.IsNotExist(err) {
		return true
	}

	log.Printf("Failed to check if issue already exists at path %s", path)
	log.Fatal(err)

	return false
}
