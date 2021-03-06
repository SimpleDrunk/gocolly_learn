package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// Course stores information about a coursera course
type Course struct {
	Title       string
	Description string
	Creator     string
	Level       string
	URL         string
	Language    string
	Commitment  string
	HowToPass   string
	Rating      string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
	// Visit only domains: coursera.org, www.coursera.org
	// colly.AllowedDomains("coursera.org", "www.coursera.org"),

	// Cache responses to prevent multiple download of pages
	// even if the collector is restarted
	// colly.CacheDir("./coursera_cache"),
	)

	// Create another collector to scrape course details
	detailCollector := c.Clone()

	courses := make([]Course, 0, 200)

	// // On every a element which has href attribute call callback
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	// If attribute class is this long string return from callback
	// 	// As this a is irrelevant

	// 	if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
	// 		return
	// 	}
	// 	link := e.Attr("href")
	// 	// If link start with browse or includes either signup or login return from callback
	// 	if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
	// 		return
	// 	}
	// 	// start scaping the page under the link found
	// 	e.Request.Visit(link)
	// })

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// On every a HTML element which has data-click-value attribute call callback
	c.OnHTML("a[data-click-value]", func(e *colly.HTMLElement) {
		// If attribute class is this long string return from callback
		// As this a is irrelevant

		if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
			return
		}
		linkValue := e.Attr("data-click-value")
		type hrefStruct struct {
			Href string `json:"href"`
		}
		var hs hrefStruct
		json.Unmarshal([]byte(linkValue), &hs)
		link := hs.Href
		if strings.Index(link, "/learn/") != -1 {
			log.Println("link", e.Request.URL.String()+link)
			detailCollector.Visit(e.Request.URL.String() + link)
		}
	})

	// Extract details of the course
	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
		log.Println("Course found", e.Request.URL)
		title := e.ChildText(".course-title")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		course := Course{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: e.ChildText("div.content"),
			Creator:     e.ChildText("div.creator-names > span"),
		}
		// Iterate over rows of the table which contains different information
		// about the course
		e.ForEach("div.ProductGlance", func(_ int, el *colly.HTMLElement) {
			switch el.ChildText("td:first-child") {
			case "Language":
				course.Language = el.ChildText("div:nth-child(5) div:nth-child(2) div:nth-child(1)")
			case "Level":
				course.Level = el.ChildText("div:nth-child(2)")
			case "Commitment":
				course.Commitment = el.ChildText("div:nth-child(2)")
			case "How To Pass":
				course.HowToPass = el.ChildText("div:nth-child(2)")
			case "User Ratings":
				course.Rating = el.ChildText("div:nth-child(2) div:nth-of-type(2)")
			}
		})
		courses = append(courses, course)
	})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://www.coursera.org")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(courses)
}
