import json
from typing import List

from models.ingredient import Ingredient


class DetailedCocktail:
    def __init__(
        self,
        name: str,
        ingredients: List[Ingredient],
        instructions: List[str],
        imagePath: str,
        glass: str,
        is_alcoholic: bool,
    ):
        self.name = name
        self.ingredients = ingredients
        self.instructions = instructions
        self.image = imagePath
        self.glass = glass
        self.is_alcoholic = is_alcoholic

    def to_json(self):
        return json.dumps(self.__dict__)
