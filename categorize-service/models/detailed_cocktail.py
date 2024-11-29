import json
from typing import List, Optional, Union
from datetime import datetime
from pydantic import BaseModel, Field, ConfigDict, field_validator

from uuid import UUID
from dateutil.parser import isoparse  # Recommended for robust datetime parsing

from utils import is_json


class IngredientItem(BaseModel):
    name: str
    amount: str

class IngredientDTO(BaseModel):
    ingredients: List[IngredientItem]

class InstructionsDTO(BaseModel):
    instructions: List[str]

class DetailedCocktail(BaseModel):
    id: UUID
    name: str
    ingredients: IngredientDTO
    instructions: Union[InstructionsDTO,str]
    image: str
    glass: str
    embedding: Optional[List[float]] = Field(default=None)
    is_alcoholic: bool
    created_at: datetime

    model_config = ConfigDict(
        populate_by_name=True,
        arbitrary_types_allowed=True
    )
    @field_validator('created_at', mode='before')
    @classmethod
    def parse_created_at(cls, v):
        if isinstance(v, str):
            return isoparse(v)
        return v

    @field_validator('ingredients', mode='before')
    @classmethod
    def parse_ingredients(cls, v):
        if isinstance(v, str):
            return json.loads(v)
        return v

    @field_validator('instructions', mode='before')
    @classmethod
    def parse_instructions(cls, v):

        if isinstance(v, str):
            if is_json(v):
                parsed = json.loads(v)
                return InstructionsDTO(instructions=parsed['instructions'])
        return v