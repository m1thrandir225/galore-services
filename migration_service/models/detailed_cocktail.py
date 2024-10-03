from dataclasses import dataclass
from models.ingredient import IngredientDto

@dataclass
class DetailedCocktail:
    id: str
    ingredients: IngredientDto
    instructions: str
    image: str
    glass: str
    isAlcoholic: bool
    name: str