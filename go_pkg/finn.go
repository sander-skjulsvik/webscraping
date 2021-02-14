package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var FINN string = "https://www.finn.no"
var FINN_REALESTATE_INDEX string = FINN + "/realestate/homes/search.html"

type Realest struct {
	Title, Address, URL, DateTime string
	ID, Price                     int
	Info                          map[string]string
	//Active                        bool
	Updates                       map[string]Realest // datetime string
}

//type RealestKey struct {
//	Title, Address string
//	ID             int
//}

func (left Realest) RightUpdates(right Realest) (Realest, bool) {
	// If diff keep Left data
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
	//if left.DateTime != right.DateTime {
	//	updates.DateTime = right.DateTime
	//	isUpdated = true
	//}
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


	if err := c.Visit(link); err != nil {
		log.Printf("Link: %s, failed. err: %s\n", link, err)
	}

	return realestCards
}

func getIndexPages_(link string, startPage int) []*string {
	indexPages := []*string{}
	var err error
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
			if pageNumber, err = strconv.Atoi(pageNumberStr); err != nil {
				log.Printf(" getIndexPages_ failed to find page number in str: %s, from href: %s \n", pageNumberStr, href)
			}
		}
		if pageNumber > startPage {
			indexPages = append(indexPages, &href)
			n++
		}
	})
	if err := c.Visit(link); err != nil {
		log.Printf("getIndexPages_: failed to visit link: %s. err: %s \n", link, err)
	}

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
	// DateTime
	r.DateTime = time.Now().String()
	// Updates
	r.Updates = make(map[string]Realest)
	// ID
	idRe, e := regexp.Compile(`(\d{8,9})$`)
	logIfErr(e, "")
	r.ID, e = strconv.Atoi(idRe.FindString(link))
	logIfErr(e, "r.ID = strconv.Atoi(idRe.FindString(link)), failed on link: "+link)


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
		// Adding data to key
		r.Info[cleanKeysForMongoDb(s.Text())] = s.Next().Text()
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

func addUpdateRealest(old Realest, new Realest) {
	if updates, isDiff := old.RightUpdates(new); isDiff {
		old.Updates[new.DateTime] = updates
	}
}

func UpdateFinnDB() {
	// Find all locations
	log.Println("Finding all locations")
	locationCodes := getAllLocations(FINN_REALESTATE_INDEX)
	locationString := "location="

	// Find all pages per location
	log.Println("Finding all pages per location. #LocationCodes:", len(locationCodes))
	var link string
	pages := []*string{}
	for location, locationCode := range locationCodes {
		fmt.Print("Location:", location)
		link = FINN_REALESTATE_INDEX + "?" + locationString + locationCode
		pages = append(pages, getIndexPages(link)...)
		fmt.Printf("(pages: %d)\n", len(pages))
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
	realestates := make(map[int]*Realest) // id -> realestate
	var currentRealestate *Realest

	for ind, link := range realestateLinks {
		fmt.Printf("realestateLinks: %d of %d\n", ind, len(realestateLinks))
		currentRealestate = getRealestateData(*link)
		if val, isIn := realestates[currentRealestate.ID]; isIn {
			// If finn id is already added look for differences. If diff add to realest.updates.
			log.Printf("Found duplicate in Find data from all listings. ID: %d \n ", currentRealestate.ID)
			addUpdateRealest(*val, *currentRealestate)
			realestates[currentRealestate.ID] = val
		} else {
			// else just add
			realestates[currentRealestate.ID] = currentRealestate
		}
		//if ind > 200 {break}
	}

	// Store data
	log.Println("Storing data. #Realestates:", len(realestates))
	collection := getFinnRealestateCollection()
	// insert new,
	// Update existing
	UpdateManyRealestate(collection, realestates)
	// Once not in the new listings mark as active: False?
}
