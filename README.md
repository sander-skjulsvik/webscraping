# Webscraping: Finn realestate

two packages, one in python using bf4 but not maintained, and a go pkg which is tried and works although is not formally tested.

## go_pkg

This is a go program for scraping realestate data of finn. The pro gram takes currently no input to run, so it just doing the same thing every time. Which is reading all realestae on finn (not all, but almost) and updates the db.

The pkg is using goquery and colly for scraping the finn pages.'

## How does it work.

In finn.go the main logic is written, and the main data flow is controlled by the function UpdateFinnDb.

TLDR: Finn realestate index -> Location links -> Pages links -> Realestate listing links -> Realestate data -> Checked for uniqueness of new once -> Checked for uniqueness in db -> Inserted to db.

1. It first gets a index page with a search query url. As currently implemented the search has no filter and therefore gives all realestates. 
2. From the index page it gets links to all top level locations. This is bc. there is a limit to how many listings per query. If a search is done without any specification only 50 pages will be avalible. By specifing location we get more pages (assumed all).
3. Whit the location links the program can get all page links. 
4. From all page links the program can get all listing links.
5. From all listing links we can get all realestate data.
6. The data is then put into a map with keys: finnkode, and value a instance of Realestate. If a key exists multiple times they are joined where the data from the newest one is put in the oldest one. The data is compared and only the difference is added to the old.updates.
7. Then the realestate instances are checked against the db. It checks if the realestate is already in the db. If so it checks for differences. If there is a difference the difference will be added to the realestate updates map, with the time as key. If the realestate is not in the db it is added. 


## To run:

The main function is in runHere.go. It also assumes a mongo db @ mongodb://localhost:27017.

## Finn.go

The finn.go file cont

## DB.

The db configuration and functions are located in db.go. There you can specify the mongodb uri. 

Fields: 

```go
type Realest struct {
	Title, Address, URL, DateTime string
	ID, Price                     int
	Info                          map[string]string
	Updates                       map[string]Realest // datetime string
	//Active                        bool
}
```

