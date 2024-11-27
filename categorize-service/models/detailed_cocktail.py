from typing import List
from datetime import datetime
from pydantic import BaseModel, Field

from models.ingredient_dto import IngredientDTO
from models.instructions_dto import InstructionsDTO
from uuid import UUID

class DetailedCocktail(BaseModel):
    id: UUID
    name: str
    ingredients: IngredientDTO
    instructions: InstructionsDTO
    image: str
    glass: str
    embedding: List[float] = Field()
    is_alcoholic: bool
    created_at: datetime

    class Config:
        allow_population_by_field_name = True