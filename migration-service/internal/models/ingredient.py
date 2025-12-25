from pydantic import BaseModel


class IngredientData(BaseModel):
    name: str
    amount: str


class IngredientDto(BaseModel):
    ingredients: list[IngredientData]

