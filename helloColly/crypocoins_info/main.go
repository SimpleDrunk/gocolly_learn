package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".cmc-view-all-coins tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText("td:nth-of-type(2)"),
			e.ChildText("td:nth-of-type(3)"),
			e.ChildText("td:nth-of-type(5)"),
			e.ChildText("td:nth-of-type(7)"),
			e.ChildText("td:nth-of-type(4)"),
			e.ChildText("td:nth-of-type(8)"),
			e.ChildText("td:nth-of-type(9)"),
			e.ChildText("td:nth-of-type(10)"),
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
