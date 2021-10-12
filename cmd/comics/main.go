package main

import (
	"flag"
	"log"
	"os"

	"github.com/antoinerey/comics/internal/series"
)

var url string
var root string

func init() {
	flag.StringVar(&root, "root", "dist", "The path to the library root")
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
		issue.Parse().Download(series.GetDirectory(root))
	}
}
