import json
import requests

from models.ingredient import Ingredient
from utils import format_ingredients, ingredients_json

cocktails_url = "https://www.thecocktaildb.com/api/json/v1/1/filter.php?c=Cocktail"


cocktail_request = requests.get(cocktails_url)

data = cocktail_request.json()["drinks"]

single_cocktail = requests.get(
    "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i=178318"
)

c_data = single_cocktail.json()["drinks"][0]

ingredients = format_ingredients(c_data)

# for cocktail in data:
# cocktail_id = cocktail["idDrink"]
# print(cocktail_id)
# single_request = requests.get(
#    "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i=" + cocktail_id
# )

# cocktail_data = single_request.json()["drinks"][0]
