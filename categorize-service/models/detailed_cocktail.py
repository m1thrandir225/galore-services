from pydantic import BaseModel

from models.ingredient_dto import IngredientDTO


class DetailedCocktail(BaseModel):
    id: str
    ingredients: IngredientDTO
    instructions: str
    image: str
    glass: str
    isAlcoholic: bool
    name: str