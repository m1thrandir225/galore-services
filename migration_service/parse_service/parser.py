import json
import time
from os.path import isfile
from typing import Dict, List
import requests
from models.detailed_cocktail import DetailedCocktail

class Parser:
    def __init__(self):
        self.url: str = (
            "https://www.thecocktaildb.com/api/json/v1/1/filter.php?c=Cocktail"
        )
        self.single_cocktail_url = (
            "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i="
        )

    def get_cocktails(self) -> List:
        response = requests.get(self.url)

        json_data: List = response.json()["drinks"]  # the list of drinks

        return json_data

    def parse_cocktails(self) -> List[DetailedCocktail]:
        cocktails: List[DetailedCocktail] = []

        data = self.get_cocktails()
        for item in data:
            drink_id = item["idDrink"]
            headers = {'Accept': 'application/json'}
            response = requests.get("{}{}".format(self.single_cocktail_url, drink_id), headers=headers)

            single_cocktail_data = response.json()["drinks"][0]
            detailed_cocktail = DetailedCocktail(single_cocktail_data)

            cocktails.append(detailed_cocktail)
            time.sleep(0.5)

        return cocktails
