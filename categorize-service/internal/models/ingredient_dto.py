from pydantic import BaseModel

from .ingredient import Ingredient


class IngredientDTO(BaseModel):
    ingredients: list[Ingredient]
