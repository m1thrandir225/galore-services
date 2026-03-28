from pydantic import BaseModel
from .ingredient import IngredientDto


class DetailedCocktail(BaseModel):
    id: str
    ingredients: IngredientDto
    instructions: str
    image: str
    glass: str
    isAlcoholic: bool
    name: str

