import json
from typing import Dict
import utils
from utils import download_image

class DetailedCocktail:
    def __init__(self, json_data: Dict[str, str]):
        self.id = json_data["idDrink"]
        self.name = json_data["strDrink"]
        self.ingredients = utils.ingredients_json(json_data)
        self.instructions = json_data["strInstructions"]
        self.image = utils.download_image(json_data["strDrinkThumb"])
        self.glass = json_data["strGlass"]
        self.is_alcoholic = utils.has_alc(json_data["strAlcoholic"])

    def to_dict(self):
        return {
            "name": self.name,
            "ingredients": self.ingredients,
            "instructions": self.instructions,
            "image": self.image,
            "glass": self.glass,
            "alcoholic": self.is_alcoholic
        }
    def to_json(self):
        return json.dumps(self.to_dict())
