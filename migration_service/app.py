import requests

cocktails_url = "https://www.thecocktaildb.com/api/json/v1/1/filter.php?c=Cocktail"


def download_image(url: str):
    print(url)


cocktail_request = requests.get(cocktails_url)

data = cocktail_request.json()["drinks"]

for cocktail in data:
    cocktail_id = cocktail["idDrink"]
    print(cocktail_id)
    # single_request = requests.get(
    #    "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i=" + cocktail_id
    # )

    # cocktail_data = single_request.json()["drinks"][0]
