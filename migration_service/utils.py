import os
from typing import Dict, List
from models.ingredient import Ingredient
import requests, json

def has_alc(alcoholic):
    if alcoholic == "Alcoholic":
        return  True
    else:
        return False

def download_image(url: str):
    response = requests.get(url)

    if response.status_code != 200:
        exit()
    filename = url.split("/")[-1]

    folderPath = "./temp"

    os.makedirs(os.path.dirname(folderPath), 511, True)

    filePath = folderPath + "/" + filename
    with open(filePath, "wb") as file:
        file.write(response.content)
    return filePath


def format_ingredients(json_data: Dict[str, str]) -> List[Ingredient]:
    ingredients = []
    ingredient_name_key = "strIngredient"
    ingredient_amount_key = "strMeasure"

    for i in range(1, 15, 1):
        if json_data["{}{}".format(ingredient_name_key, i)] is None:
            break

        name = json_data["{}{}".format(ingredient_name_key, i)]
        amount = json_data["{}{}".format(ingredient_amount_key, i)]

        ingredient = Ingredient(name, amount)
        ingredients.append(ingredient)

    return ingredients


def ingredients_json(json_data: Dict[str, str]) -> str:
    ingredients = format_ingredients(json_data)
    list_json = json.dumps([ingredient.to_dict() for ingredient in ingredients])

    return json.dumps(list_json)


#download_image("https://www.thecocktaildb.com/images/media/drink/wysqut1461867176.jpg")
