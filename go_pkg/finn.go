package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var FINN string = "https://www.finn.no"
var FINN_REALESTATE_INDEX string = FINN + "/realestate/homes/search.html"

type Realest struct {
	Title, Address, URL string
	ID, Price           int
	Info                map[string]string
}

func realests2csv(reals []Realest) {

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

	c.Visit("https://www.finn.no/realestate/homes/search.html?sort=PUBLISHED_DESC")

	return realestCards
}

func getIndexPages(link string, startPage int) []*string {
	indexPages := []*string{}
	re, _ := regexp.Compile(`\d{1,2}`)
	n := 0

	c := colly.NewCollector(
		colly.AllowedDomains("finn.no", "www.finn.no"),
		//colly.Async(true),
	)

	c.OnHTML("[class='pagination__page button button--pill']", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		x := re.FindString(href)
		pageNumber, _ := strconv.Atoi(x)
		if (pageNumber >= startPage) || (x == "") {
			indexPages = append(indexPages, &href)
			n++
		}
	})
	c.Visit(link)

	if n > 0 {
		newLink := FINN_REALESTATE_INDEX + *indexPages[n-1]
		var newStartPage = startPage + n
		return append(indexPages, getIndexPages(newLink, newStartPage)...)
	}
	return indexPages
}

func getRealestateData(link string) *Realest {
	var r Realest
	var price string
	doc, err := goquery.NewDocument(link)
	if err != nil {
		log.Println("Err on", link, "err: ", err)
	}
	mainArea := doc.Find("div.grid")
	title := mainArea.Find("h1")
	r.Title = title.Text()
	r.Address = title.Next().Text()
	r.URL = link
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

	//pages := getIndexPages(FINN_REALESTATE_INDEX, 1)
	////fmt.Println("pages: ")
	////PrintStingArr(pages)
	//
	//var cards []*string
	//for _, link := range pages {
	//	x := getRealestateCardUrls(*link)
	//	cards = append(cards, x...)
	//
	//}
	//fmt.Println("cards")
	//PrintStingArr(cards)
	//
	//r := getRealestateData("https://www.finn.no/realestate/homes/ad.html?finnkode=205519251")
	//log.Println("Title:", r.Title)
	//log.Println("Addr:", r.Address)
	//log.Println("Price :", r.Price)
	//log.Println("Info:", r.Info)

	fmt.Println(getAllLocations("https://www.finn.no/realestate/homes/search.html?sort=PUBLISHED_DESC"))

}
