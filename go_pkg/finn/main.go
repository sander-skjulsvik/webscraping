package finn

import (
	"fmt"
	"log"

	db "github.com/sander-skjulsvik/webscraping/go_pkg/db"
)

var FINN string = "https://www.finn.no"
var FINN_REALESTATE_INDEX string = FINN + "/realestate/homes/search.html"

func getIndexes(startURL string, indexes chan<- string) {
	// Find locations
	log.Println("Finding locations")
	locationCodes := GetAllLocations(startURL)
	locationString := "location="

	// Find pages per location
	log.Println("Finding pages per location. #LocationCodes:", len(locationCodes))
	for location, locationCode := range locationCodes {
		fmt.Println("Location:", location)
		indexPages := GetIndexPages(FINN_REALESTATE_INDEX + "?" + locationString + locationCode)
		for _, indexLink := range indexPages {
			indexes <- *indexLink
		}
		log.Printf("#indexes +%d: %d\n", len(indexPages), len(indexes))
	}
}

func getRealestates(indexes <-chan string, realestateLinks chan<- string) {
	// Find listings per page
	for page := range indexes {
		cardUrls := getRealestateCardUrls(FINN_REALESTATE_INDEX + page)
		for _, realestateLink := range cardUrls {
			realestateLinks <- *realestateLink
		}
		log.Printf("#realestateLinks +%d: %d\n", len(cardUrls), len(realestateLinks))
		log.Printf("#indexes -1: %d\n", len(realestateLinks))
	}
}

func processRealestateLinks(realestateLinks <-chan string, realestateData chan<- db.Realestate) {
	// Find data from listings
	log.Println("Finding data for listings. #RealestateLinks:", len(realestateLinks))
	processed := make(map[string]bool) // url -> true

	for link := range realestateLinks {
		fmt.Printf("realestateLinks: %d\n", len(realestateLinks))
		if _, isIn := processed[link]; !isIn {
			realestateData <- GetRealestateData(link)
			processed[link] = true
		} else {
			log.Printf("Skipping duplicate link: %s\n", link)
		}
		log.Printf("#realestateLinks -1: %d\n", len(realestateLinks))
		log.Printf("#realestateData +1: %d\n", len(realestateData))
	}
}

func addRealestateDataToDB(realestateData <-chan db.Realestate) {
	// Store data
	log.Println("Storing data")
	collection := db.GetFinnRealestateCollection()
	// insert new,
	// Update existing
	for realestate := range realestateData {
		db.UpdateRealestate(collection, realestate)
		log.Printf("#realestateData: %d\n", len(realestateData))
	}
	// Once not in the new listings mark as active: False?

}

func UpdateFinnDB() {
	// Prep
	realestateIndexLinks := make(chan string, 100)
	realestateLinks := make(chan string, 1000)
	realestateData := make(chan db.Realestate, 10000)

	// go produce: Get index pages
	go getIndexes(FINN_REALESTATE_INDEX, realestateIndexLinks)

	// go consume produce: get realestate links from index pages and add to process que
	go getRealestates(realestateIndexLinks, realestateLinks)

	// go consume and produce: process links add to que of structs with data
	go processRealestateLinks(realestateLinks, realestateData)

	// go consume: add to db
	go addRealestateDataToDB(realestateData)

}
