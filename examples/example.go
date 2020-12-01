package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"

	"github.com/MicahParks/collymongo"
)

func main() {

	// Create the collector.
	c := colly.NewCollector()

	// Create the MongoDB storage backend.
	storage := &collymongo.CollyMongo{Uri: "mongodb://localhost:27017"}

	// Set the storage backend.
	if err := c.SetStorage(storage); err != nil {
		log.Fatalln(err)
	}

	// Find and visit all links.
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if err := e.Request.Visit(e.Attr("href")); err != nil {

			// If the link has been visited before or if the URL is empty, then skip it.
			if err.Error() != "URL already visited" && err.Error() != "Missing URL" {
				log.Fatalln(err)
			}
		}
	})

	// State what URL the scraper is on.
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: " + r.URL.String())
	})

	// Start the scraper on the Go Colly website.
	if err := c.Visit("http://go-colly.org/"); err != nil {
		log.Fatalln(err)
	}
}
