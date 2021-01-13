from bs4 import BeautifulSoup
from urllib.parse import urlparse
import requests

from util import *

from IPython import embed


def get_soup(url):
    page = requests.get(url)
    return BeautifulSoup(page.content, 'html.parser')


def find_cards(soup):
    cards = {card.h2.get_text(): card.h2.a["href"] for card in soup.find_all(
        "article", {"class": "ads__unit"})}
    return cards


def find_all_pages(soup, start_side: int, base_url):
    pages = []
    current_page = start_side
    for page in soup.find_all("a", {"pagination__page button button--pill"}):
        current_page = extract_int(page["aria-label"])
        if current_page > start_side:
            pages.append(page)
    if len(pages) == 0:
        return []
    new_soup = get_soup(f"{base_url}{pages[-1]['href']}")
    return pages + find_all_pages(
        new_soup, current_page, base_url)


def get_data(soup):
    kw = "definition-list definition-list--cols1to2"
    dl_table = soup.find("dl", {"class": kw})
    return {dt.get_text(): dt.find_next("dd").get_text()
            for dt in dl_table.find_all("dt")}


if __name__ == "__main__":
    # get index
    base_url = "https://www.finn.no"
    search_url = f"{base_url}/realestate/homes/search.html"
    index_url = f"{search_url}?location=0.20061&sort=PUBLISHED_DESC"
    soup = get_soup(index_url)

    # get all pages
    pages = find_all_pages(
        soup, start_side=0, base_url="https://www.finn.no/realestate/homes/search.html")

    # get all cards

    cards = {}
    for page in pages:
        cards.update(
            find_cards(get_soup(search_url + page["href"]))
        )
    # making all links absolute,because not all was.
    for key, val in cards.items():
        if "http" not in val:
            cards[key] = "https://www.finn.no" + val

    embed()
    # get data from cards
