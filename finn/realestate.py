from typing import List
import pandas as pd
from datetime import date


from . import util

from IPython import embed


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
    new_soup = util.get_soup(f"{base_url}{pages[-1]}")
    return pages + get_all_page_links(
        new_soup, current_page, base_url)


def get_data(link):
    soup = util.get_soup(link).find(
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
    data.update({dt.get_text(strip=True): dt.find_next("dd").get_text(strip=True)
                 for dt in about_table.find_all("dt")})

    # Definition list
    # data["price_info"] = {}
    dls = panel.find_all("dl")
    for dl in dls:
        for dt in dl.find_all("dt"):
            data[dt.get_text(strip=True)] = dt.find_next(
                "dd").get_text(strip=True)

    # Cleaning of data
    if "Bruksareal" in data:
        data["Bruksareal"] = util.extract_int(data['Bruksareal'][:-1])
    if "Totalpris" in data:
        data["Totalpris"] = util.extract_int(data["Totalpris"])

    return data


# get index
BASE_URL = "https://www.finn.no"
BASE_SEARCH_URL = f"{BASE_URL}/realestate/homes/search.html"
SEARCH_STR: str = f"sort=PUBLISHED_DESC"


def main(location_code: str = None, out: str = f"out/{date.today()}_realestate", verbose=True) -> pd.DataFrame:

    base_url = BASE_URL
    search_url = BASE_SEARCH_URL
    if location_code:
        index_url = f"{BASE_SEARCH_URL}?location={location_code}&sort=PUBLISHED_DESC"
        out += "_" + location_code
    else:
        index_url = f"{BASE_SEARCH_URL}?{SEARCH_STR}"

    if verbose:
        print(f"index: {index_url}")

    # get index
    index_soup = util.get_soup(index_url)

    # get all pages
    page_links = get_all_page_links(
        index_soup, start_side=0, base_url=search_url)

    if verbose:
        print(f"Number of pages: {len(page_links)}")

    # get all cards
    card_links = list()
    for ind, page_link in enumerate(page_links):
        r = get_card_links(util.get_soup(search_url + page_link))
        for card_link in r:
            if "http" in card_link:
                card_links.append(card_link)
            # making all links absolute,because not all was.
            else:
                card_links.append("https://www.finn.no" + card_link)
            # embed()

    if verbose:
        print(f"Number of cards: {len(card_links)}")

    # get data from cards
    cards_data = []
    N = len(card_links)
    N_failed = 0
    for i, card_link in enumerate(card_links):
        if "nybyg" in card_link:
            msg = f"Nybygg: url: {card_link}"
            with open("log/all_realestate_oslo_log.txt", "a+") as fp:
                fp.write(msg + "\n")
            N_failed += 1
            continue
        try:
            cards_data.append(
                get_data(card_link)
            )
        except AttributeError as e:
            msg = f"url: {card_link}\n{e}"
            with open("log/all_realestate_oslo_log.txt", "a+") as fp:
                fp.write(msg + "\n")
            print(msg)
            N_failed += 1
        print(f"{(i/N)*100}. {i} out of {N}, {N_failed=}")

    # df
    df = pd.DataFrame.from_dict(cards_data)

    # write csv
    # updating the base one
    df.to_csv("out/realestate.csv")
    # wrting one new with specified name, default is with date.
    df.to_csv(out + ".csv")
    return df
