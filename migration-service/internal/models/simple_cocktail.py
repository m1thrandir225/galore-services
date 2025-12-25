from pydantic import BaseModel


class SimpleCocktail(BaseModel):
    strDrink: str
    strDrinkThumb: str
    idDrink: str
