from typing import List

from models.simple_cocktail import SimpleCocktail


class Cocktails:
    def __init__(self, cocktails: List[SimpleCocktail]):
        self.cocktails: List[SimpleCocktail] = cocktails
