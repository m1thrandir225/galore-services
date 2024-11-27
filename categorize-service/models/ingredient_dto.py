from typing import List

from pydantic import BaseModel

from models.ingredient import Ingredient


class IngredientDTO(BaseModel):
    ingredients: List[Ingredient]