import os
from typing import Dict, List
from models.ingredient import IngredientData, IngredientDto
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

    folder_path = "./temp"

    os.makedirs(os.path.dirname(folder_path), 511, True)

    file_path = folder_path + "/" + filename
    with open(file_path, "wb") as file:
        file.write(response.content)
    return file_path


def format_ingredients(json_data: Dict[str, str]) -> IngredientDto:
    ingredients = []
    ingredient_name_key = "strIngredient"
    ingredient_amount_key = "strMeasure"

    for i in range(1, 15, 1):
        if json_data["{}{}".format(ingredient_name_key, i)] is None:
            break

        name = json_data["{}{}".format(ingredient_name_key, i)]
        amount = json_data["{}{}".format(ingredient_amount_key, i)]

        ingredient = IngredientData(name, amount)
        ingredients.append(ingredient)

    return IngredientDto(ingredients)
