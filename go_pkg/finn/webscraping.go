package finn

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sander-skjulsvik/webscraping/go_pkg/db"
	"github.com/sander-skjulsvik/webscraping/go_pkg/utils"
)

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

func GetIndexPages(link string) []*string {
	return getIndexPages_(link, 0)
}

func GetRealestateData(link string) db.Realestate {

	var r db.Realestate
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
	// ID
	idRe, e := regexp.Compile(`(\d{8,9})$`)
	utils.LogIfErr(e, "")
	r.ID, e = strconv.Atoi(idRe.FindString(link))
	utils.LogIfErr(e, "r.ID = strconv.Atoi(idRe.FindString(link)), failed on link: "+link)

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
	r.Price, err = strconv.Atoi(price)
	if err != nil {
		r.Price = -1
		log.Printf("Realestate: %s, could not convert price to int: %s. Setting price to -1\n", r.Address, price)
	}
	// get all dls/Info
	r.Info = make(map[string]string)
	mainArea.Find("dt").Each(func(i int, s *goquery.Selection) {
		// Adding data to key
		r.Info[utils.CleanKeysForMongoDb(s.Text())] = s.Next().Text()
	})
	return r
}

func GetAllLocations(link string) map[string]string {
	locations := make(map[string]string)
	doc, err := goquery.NewDocument(link)
	removeParenthesis, _ := regexp.Compile("\\([^)]*\\)")
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
		key := removeParenthesis.ReplaceAllString(s.Text(), "")
		v, _ := s.Find("label").Attr("for")
		val := extractFloat.FindString(v)
		locations[key] = val

	})
	//fmt.Println(h3.Text())

	return locations
}
