from bs4 import BeautifulSoup
import re
import requests


def extract_int(s) -> int:
    return int("".join(digit for digit in s if digit.isdigit()))


def isfloatable(s: str) -> bool:
    """
    Somethimes is easier to ask for forgivness
    """
    try:
        return isinstance(float(s), float)
    except:
        return False


def get_soup(url):
    page = requests.get(url)
    return BeautifulSoup(page.content, 'html.parser')


def get_all_locations(url: str = "https://www.finn.no/realestate/homes/search.html?sort=PUBLISHED_DESC"):
    soup = get_soup(url)
    h3 = [h3 for h3 in soup.find_all(
        "h3", {"class": "u-t5 u-mt32"}) if h3.get_text() == "Område"][0]
    return {
        re.sub(r' \([^)]*\)', "", li.get_text()): str(li.input)[20:27]
        for li in h3.find_next("ul").find_all("li")
    }
