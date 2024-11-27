from pydantic import BaseModel

class Ingredient(BaseModel):
    name: str
    description: str