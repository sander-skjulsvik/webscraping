
# %%

import argparse
from finn import realestate, util
from datetime import date


LOCATIONS = {
    "agder": "0.22042",
    "innlandet": "0.22034",
    "møre og romsdal": "0.22034",
    "Nordland": "0.20018",
    "oslo": "0.20061",
    "rogaland": "0.20012",

}


def update_finn_realestate(location: str = None, out: str = None, verbose=False):
    """
    Location oslo: 0.20061
    """
    # Handle default args
    if location and out:
        # if locaton is a location name and not a code gather the codes with get all locations
        if not util.isfloatable(location):
            location_map = util.get_all_locations()
            try:
                location = location_map[location]
            except KeyError as e:
                print(
                    f"Given location name is is not known. Known locations {location_map}")
                exit()
        realestate.main(location, out)
    elif location:
        realestate.main(location=location, verbose=True)

    if verbose:
        print("choosing all known locations")
    location_map = util.get_all_locations()
    if not out:
        out = "out/"
    out = f"{date.today()}_" + out + "_{}"
    for ind, location in enumerate(location_map):
        print(
            f"Starting location: {location}, ({ind}/{len(location_map)}={ind/len(location_map)})")
        realestate.main(location_map[location], out=out.format(location))


# %%
funcs = {
    "update_finn_realestate": update_finn_realestate,

}
# %%
if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("func", choices=list(funcs.keys()),
                        help=f"add one of thise functions {funcs.keys()}", nargs='?', default="update_finn_realestate")
    parser.add_argument("--args", nargs='+',
                        help=f"For update_finn_realestate, locations are {util.get_all_locations()}")

# %%
    args = parser.parse_args()
# %%
    if args.args is None:
        funcs[args.func]()
    else:
        funcs[args.func](*args.args)
    # funcs[args.func](*args.args)
