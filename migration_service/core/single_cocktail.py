import requests


class SingleCocktail:
    def __init__(
        self,
        name: str,
        instructions: str,
        ingredients: list[str],
        glass: str,
        is_alcoholic: bool,
        idDrink: str,
    ):
        self.name = name
        self.instructions = instructions
        self.ingredients = ingredients
        self.glass = glass
        self.is_alcoholic = is_alcoholic
        self.idDrink = idDrink

    def get_cocktail():
        r = requests.get(url)
        data = r.json()
