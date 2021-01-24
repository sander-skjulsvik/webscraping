# Webscraping repo

## Finn

Web scraping realestate data from finn.no.

The logic of the scraping is in the all_realestate_finn_oslo.py script. It is using Beautifulsoup for the scraping, and pandas for converting from python dictionary to csv. 

The script uses a url to the finn page defined in the global variables shown underneath.

```Python
BASE_URL = "https://www.finn.no"
BASE_SEARCH_URL = f"{BASE_URL}/realestate/homes/search.html"
SEARCH_STR: str = f"location=0.20061&sort=PUBLISHED_DESC"
```
To edit the search changed the SEARCH_STR with the url to the index of the search result. Then the script will find all pages and all cards on those pages according to the search. Then it will gather the data from all listings.


### Updating the data

```shell
$ python3 interface.py update_finn_realestate
```

### Data

Realestate data saved in out/

