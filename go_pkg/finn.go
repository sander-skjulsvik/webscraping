package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var FINN string = "https://www.finn.no"
var FINN_REALESTATE_INDEX string = FINN + "/realestate/homes/search.html"

type Realest struct {
	Title, Address, URL, Date string
	ID, Price                 int
	Info                      map[string]string
}

func (left Realest) RightUpdates(right Realest) (Realest, bool) {
	// If diff keep left data
	updates := Realest{}
	isUpdated := false
	if left.Title != right.Title {
		updates.Title = right.Title
		isUpdated = true
	}
	if left.Address != right.Address {
		updates.Address = right.Address
		isUpdated = true
	}
	if left.URL != right.URL {
		updates.URL = right.URL
		isUpdated = true
	}
	if left.Date != right.Date {
		updates.Date = right.Date
		isUpdated = true
	}
	if left.ID != right.ID {
		updates.ID = right.ID
		isUpdated = true
	}
	if left.Price != right.Price {
		updates.Price = right.Price
		isUpdated = true
	}
	for rightKey, rightVal := range right.Info {
		leftVal, ok := left.Info[rightKey]
		if !ok || leftVal != rightVal {
			updates.Info[rightKey] = rightVal
			isUpdated = true
		}
	}
	return updates, isUpdated
}
func (left Realest) LeftUpdates(right Realest) (Realest, bool) {
	return right.RightUpdates(left)
}


func getRealestateCardUrls(link string) []*string {
	c := colly.NewCollector(
		colly.AllowedDomains("finn.no", "www.finn.no"),
		//colly.Async(true),
	)

	realestCards := []*string{}

	c.OnHTML(`article`, func(elm *colly.HTMLElement) {
		r := elm.ChildAttr("a", "href")
		if !strings.Contains(r, "nybygg") {
			if !strings.Contains(r, "http") {
				r = FINN + r
			}
			realestCards = append(realestCards, &r)
		}
	})

	c.Visit(link)

	return realestCards
}

func getIndexPages_(link string, startPage int) []*string {
	indexPages := []*string{}
	reFindPage, _ := regexp.Compile(`(page=)\d{1,2}`)
	reRemovePageStr, _ := regexp.Compile(`(page=)`)
	n := 0

	c := colly.NewCollector(
		colly.AllowedDomains("finn.no", "www.finn.no"),
		//colly.Async(true),
	)

	c.OnHTML("[class='pagination__page button button--pill']", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		pageNumber := 1
		if strings.Contains(href, "page") {
			pageNumberStr := reRemovePageStr.ReplaceAllString(reFindPage.FindString(href), "")
			pageNumber, _ = strconv.Atoi(pageNumberStr)
		}
		if pageNumber > startPage {
			indexPages = append(indexPages, &href)
			n++
		}
	})
	c.Visit(link)

	if n > 0 {
		newLink := FINN_REALESTATE_INDEX + *indexPages[n-1]
		var newStartPage = startPage + n
		return append(indexPages, getIndexPages_(newLink, newStartPage)...)
	}
	return indexPages
}

func getIndexPages(link string) []*string {
	return getIndexPages_(link, 0)
}

func getRealestateData(link string) *Realest {
	//Title, Address, URL string
	//ID, Price           int
	//Info                map[string]string
	// TODO : not keys ending with "."
	var r Realest
	var e error

	var price string
	doc, err := goquery.NewDocument(link)
	if err != nil {
		log.Println("Err on", link, "err: ", err)
	}
	mainArea := doc.Find("div.grid")
	title := mainArea.Find("h1")
	// Title
	r.Title = title.Text()
	// Address
	r.Address = title.Next().Text()
	// URL
	r.URL = link
	// Date
	// TODO : Change to only date
	r.Date = time.Now().String()
	// ID
	idRe, e := regexp.Compile(`(\d{8,9})$`)
	logIfErr(e, "")
	r.ID, e = strconv.Atoi(idRe.FindString(link))
	if logIfErr(e, "r.ID = strconv.Atoi(idRe.FindString(link)), failed on link: "+link) {
		r.ID = 0
		log.Printf("")
	}
	// Get price
	mainArea.Find("span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		msg := s.Text()
		if strings.Contains(msg, "Prisantydning") {
			price = s.Next().Text()
			price = price[:len(price)-2]
			// return false to break early.
			return false
		}
		return true
	})
	// Price
	r.Price, _ = Ascii2Int(price)

	// get all dls/Info
	r.Info = make(map[string]string)
	mainArea.Find("dt").Each(func(i int, s *goquery.Selection) {
		r.Info[s.Text()] = s.Next().Text()
	})
	return &r
}

func getAllLocations(link string) map[string]string {
	locations := make(map[string]string)
	doc, err := goquery.NewDocument(link)
	removeParanthesis, _ := regexp.Compile("\\([^)]*\\)")
	extractFloat, _ := regexp.Compile("0.\\d+")

	if err != nil {
		log.Println("Err on", link, "err: ", err)
	}
	var h3 *goquery.Selection
	doc.Find("h3.u-t5").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Text() == "Område" {
			h3 = s
			return false
		}
		return true
	})

	ul := h3.Next()
	ul.Find("li").Each(func(i int, s *goquery.Selection) {
		key := removeParanthesis.ReplaceAllString(s.Text(), "")
		v, _ := s.Find("label").Attr("for")
		val := extractFloat.FindString(v)
		locations[key] = val

	})
	//fmt.Println(h3.Text())

	return locations
}

func main() {

	// Find all locations
	log.Println("Finding all locations")
	locationCodes := getAllLocations(FINN_REALESTATE_INDEX)
	locationString := "location="

	// Find all pages per location
	log.Println("Finding all pages per location. #LocationCodes:", len(locationCodes))
	var link string
	pages := []*string{}
	for location, locationCode := range locationCodes {
		fmt.Println("Location:", location)
		link = FINN_REALESTATE_INDEX + "?" + locationString + locationCode
		pages = append(pages, getIndexPages(link)...)
	}
	// Find all listings per page
	log.Println("Finding all listings per page. #Pages:", len(pages))
	realestateLinks := []*string{}
	for ind, page := range pages {
		fmt.Printf("Page: %d of %d \n", ind, len(pages))
		realestateLinks = append(realestateLinks, getRealestateCardUrls(FINN_REALESTATE_INDEX+*page)...)
	}
	// Find data from all listings
	log.Println("Finding all data for all listings. #RealestateLinks:", len(realestateLinks))
	realestates := []interface{}{Realest{}}
	for ind, link := range realestateLinks {
		fmt.Printf("realestateLinks: %d of %d\n", ind, len(realestateLinks))
		realestates = append(realestates, getRealestateData(*link))
	}
	// Store data
	log.Println("Storing data. #Realestates:", len(realestates))
	collection := getFinnRealestateCollection()
	insertManyRealestate(collection, realestates)

}
