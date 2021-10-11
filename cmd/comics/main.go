package main

import (
	"log"
	"os"

	"github.com/antoinerey/comics/internal/series"
)

func main() {
	url := os.Args[1]

	series, err := series.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(series)
}
