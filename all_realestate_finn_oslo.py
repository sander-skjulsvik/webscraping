from typing import List
from bs4 import BeautifulSoup
from urllib.parse import urlparse
import requests
import pandas as pd

import util
from IPython import embed


def get_soup(url):
    page = requests.get(url)
    return BeautifulSoup(page.content, 'html.parser')


def get_card_links(soup):
    cards = {card.h2.a["href"] for card in soup.find_all(
        "article", {"class": "ads__unit"})}
    return cards


def get_all_page_links(soup, start_side: int, base_url) -> List[str]:
    """
    Getting links to all pages in the index

    Args:
        soup: bf4.BeautifulSoup
            index soup
        start_site: int
            Side to incude from
        base url: str
            base url for adding to href
    Return: List[str]
        list of hrefs to pages
    """
    pages = []
    current_page = start_side
    for page in soup.find_all("a", {"pagination__page button button--pill"}):
        current_page = util.extract_int(page["aria-label"])
        if current_page > start_side:
            pages.append(page["href"])
    if len(pages) == 0:
        return []
    new_soup = get_soup(f"{base_url}{pages[-1]}")
    return pages + get_all_page_links(
        new_soup, current_page, base_url)


def get_data(link):
    soup = get_soup(link).find(
        "main", {"class": "pageholder"})

    data = {
        "url": link,
    }
    # header
    header = soup.find("h1")
    data["header"] = header.get_text(strip=True)

    # location
    location = header.find_next_siblings('p')[0]
    data["location"] = location.get_text(strip=True)

    panel = location.find_next("div", {"class": "panel"})

    # price info
    price = panel.find_all("span")[1]
    data["price"] = util.extract_int(price.get_text(strip=True))

    # type of realestate
    kw = "definition-list definition-list--cols1to2"
    about_table = soup.find("dl", {"class": kw})
    data["about"] = {dt.get_text(strip=True): dt.find_next("dd").get_text(strip=True)
                     for dt in about_table.find_all("dt")}

    # Definition list
    data["price_info"] = {}
    dls = panel.find_all("dl")
    for dl in dls:
        for dt in dl.find_all("dt"):
            data["price_info"][dt.get_text(strip=True)] = dt.find_next(
                "dd").get_text(strip=True)

    return data


"""

- Get index
-> get all pages
-> get all cards
->
"""
a = [1]

if __name__ == "__main__":
    # get index
    base_url = "https://www.finn.no"
    search_url = f"{base_url}/realestate/homes/search.html"
    index_url = f"{search_url}?location=0.20061&sort=PUBLISHED_DESC"
    index_soup = get_soup(index_url)

    # get all pages
    page_links = get_all_page_links(
        index_soup, start_side=0, base_url=search_url)

    # get all cards
    card_links = set()
    for page_link in page_links:
        for card_link in get_card_links(get_soup(search_url + page_link)):
            if "http" in card_link:
                card_links.add(card_link)
            # making all links absolute,because not all was.
            else:
                card_links.add("https://www.finn.no" + card_link)
            # embed()

    # get data from cards
    cards_data = []
    N = len(card_links)
    N_failed = 0
    for i, card_link in enumerate(card_links):
        try:
            cards_data.append(
                get_data(card_link)
            )
        except AttributeError as e:
            print(f"url: {card_link}\n{e}")
            print(e)
            N_failed += 1
        print(f"{(i/N)*100}. {i} out of {N}, {N_failed=}")

    # df
    df = pd.DataFrame.from_dict(cards_data)

    # parsing

    # write csv
    df.to_csv("realestate.csv")
