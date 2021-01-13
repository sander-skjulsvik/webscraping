from typing import List
from bs4 import BeautifulSoup
from urllib.parse import urlparse
import requests

import util
from IPython import embed


def get_soup(url):
    page = requests.get(url)
    return BeautifulSoup(page.content, 'html.parser')


def get_card_links(soup):
    cards = {card.h2.get_text(): card.h2.a["href"] for card in soup.find_all(
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
    soup = get_soup("https://www.finn.no/realestate/homes/ad.html?finnkode=204438115").find(
        "main", {"class": "pageholder"})

    data = {
        "url": link,
    }
    # header
    header = soup.find("h1")
    data["header"] = header.get_text()

    # location
    location = header.find_next_siblings('p')[0]
    data["location"] = location.get_text()

    # panel info
    panel = location.find_next("div", {"class": "panel"})
    price = panel.find_all("span")[1]
    data["price"] = price.get_text()

    # type of realestate
    kw = "definition-list definition-list--cols1to2"
    about_table = soup.find("dl", {"class": kw})
    data["about"] = {dt.get_text(): dt.find_next("dd").get_text()
                     for dt in about_table.find_all("dt")}

    # Definition list
    dls = panel.find_all("dl")
    for dl in dls:
        pass

    data["definition_list"]

    return data


"""

- Get index
-> get all pages
-> get all cards
->

python, c/c++, java, Rust, go, kotlin, R, matlab, (p)SQl


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
        index_soup, start_side=0, base_url="https://www.finn.no/realestate/homes/search.html")

    # get all cards
    card_links = []
    for page_link in page_links:
        for card_link in get_card_links(get_soup(search_url + page_link)):
            if "http" in card_link:
                card_links.append(card_link)
            # making all links absolute,because not all was.
            else:
                card_links.append("https://www.finn.no" + card_link)
            # embed()

    # get data from cards
    card_data = []
    N = len(card_links)
    N_failed = 0
    for i, card_link in enumerate(card_links):
        try:
            card_data.append(
                get_data(card_link)
            )
        except AttributeError as e:
            print(f"url: {card_link}\n{e}")
            print(e)
            N_failed += 1
        # p = (i/N)*100
        print(f"{(i/N)*100}. {i} out of {N}, {N_failed=}")
