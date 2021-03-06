package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type cryptocoin struct {
	Name      string
	Symbol    string
	Price     string
	Volume    string
	Capacity  string
	Change1h  string
	Change24h string
	Change7d  string
}

func main() {
	fName := "cryptocoinmarketcap.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector()

	totalJSON := make(map[string]cryptocoin, 0)
	c.OnHTML(".cmc-view-all-coins tbody tr", func(e *colly.HTMLElement) {
		id := e.ChildText("td:nth-of-type(1)")
		c := cryptocoin{
			Name:      e.ChildText("td:nth-of-type(2)"),
			Symbol:    e.ChildText("td:nth-of-type(3)"),
			Price:     e.ChildText("td:nth-of-type(5)"),
			Volume:    e.ChildText("td:nth-of-type(7)"),
			Capacity:  e.ChildText("td:nth-of-type(4)"),
			Change1h:  e.ChildText("td:nth-of-type(8)"),
			Change24h: e.ChildText("td:nth-of-type(9)"),
			Change7d:  e.ChildText("td:nth-of-type(10)"),
		}
		totalJSON[id] = c
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	// c.Wait()
	b, err := json.Marshal(&totalJSON)
	if err != nil {
		log.Println("json marshal failed, err:", err)
		return
	}
	file.Write(b)
	log.Printf("Scraping finished, check file %q for results\n", fName)
}
