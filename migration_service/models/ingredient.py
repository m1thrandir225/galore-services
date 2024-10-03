import json
from dataclasses import dataclass
from typing import List


@dataclass
class IngredientData:
    name: str
    amount: str

@dataclass
class IngredientDto:
    ingredients: List[IngredientData]