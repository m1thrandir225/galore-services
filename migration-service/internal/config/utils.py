import os
from ..models import IngredientData, IngredientDto
import requests


def has_alc(alcoholic: str) -> bool:
    if alcoholic == "Alcoholic":
        return True
    else:
        return False


def download_image(url: str) -> str:
    response = requests.get(url)

    if response.status_code != 200:
        exit()
    filename = url.split("/")[-1]

    folder_path = "./temp"

    os.makedirs(folder_path, mode=0o777, exist_ok=True)

    file_path = folder_path + "/" + filename
    with open(file_path, "wb") as file:
        file.write(response.content)
    return file_path


def format_ingredients(json_data: dict[str, str]) -> IngredientDto:
    ingredients: list[IngredientData] = []
    ingredient_name_key = "strIngredient"
    ingredient_amount_key = "strMeasure"

    for i in range(1, 15, 1):
        if json_data["{}{}".format(ingredient_name_key, i)] is None:
            break

        name = json_data["{}{}".format(ingredient_name_key, i)]
        amount = json_data["{}{}".format(ingredient_amount_key, i)]

        ingredient = IngredientData(name=name, amount=amount)
        ingredients.append(ingredient)

    return IngredientDto(ingredients=ingredients)
