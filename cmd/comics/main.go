package main

import (
	"flag"
	"log"
	"os"

	"github.com/antoinerey/comics/internal/series"
)

var url string
var root string
var missing bool

func init() {
	flag.StringVar(&root, "root", "dist", "The path to the library root")
	flag.BoolVar(&missing, "missing", true, "Only download missing issues")
	flag.Parse()

	url = os.Args[1]
}

func main() {
	series, err := series.CreateSeries(url).Parse()
	if err != nil {
		log.Printf("Failed to parse series %s", url)
		log.Fatal(err)
	}

	for _, issue := range series.Issues {
		issue = issue.Parse()

		if missing && !issue.IsMissing(series.GetDirectory(root)) {
			log.Printf("Skipping %s. It's already been downloaded", issue.Title)
			continue
		}

		issue.Download(series.GetDirectory(root))
	}
}
