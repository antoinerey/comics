package main

import (
	"flag"
	"log"
	"os"

	"github.com/antoinerey/comics/internal/series"
)

var url string
var baseDir string
var tmpDir string
var missing bool

func init() {
	flag.StringVar(&baseDir, "baseDir", "dist", "The path to the library base directory")
	flag.StringVar(&tmpDir, "tmpDir", "/tmp", "The path to the temporary directory")
	flag.BoolVar(&missing, "missing", true, "Only download missing issues")
	flag.Parse()

	url = os.Args[len(os.Args)-1]
}

func main() {
	series, err := series.CreateSeries(url).Parse()
	if err != nil {
		log.Printf("Failed to parse series %s", url)
		log.Fatal(err)
	}

	for _, issue := range series.Issues {
		issue = issue.Parse()

		if missing && !issue.IsMissing(series.GetDirectory(baseDir)) {
			log.Printf("Skipping %s. It's already been downloaded", issue.Title)
			continue
		}

		issue.Download(series.GetDirectory(baseDir), series.GetDirectory(tmpDir))
	}
}
