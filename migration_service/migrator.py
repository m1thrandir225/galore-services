import json

import requests
import logging
from typing import List

import utils
from models.detailed_cocktail import DetailedCocktail


class Migrator:
    def __init__(self, url: str):
        self.url = url
        self.updated_cocktails_cache = set()  # Cache to store successfully updated cocktail names or IDs

    def update_cocktails(self, cocktails: List[DetailedCocktail]):
        """Updates cocktails by sending data to the given URL with a caching mechanism."""
        for cocktail in cocktails:
            # Use cocktail name or ID as a unique identifier for the cache
            if cocktail.name in self.updated_cocktails_cache:
                logging.info(f"Skipping {cocktail.name}, already updated.")
                continue

            try:
                files = {
                    "file": ("random.jpg", open(cocktail.image, "rb"), "image/jpg")
                }
                data = {
                    "name": cocktail.name,
                    "ingredients": json.dumps(json.dumps({
                        "ingredients": [ingredient.__dict__ for ingredient in cocktail.ingredients]
                    })),
                    "instructions": cocktail.instructions,
                    "glass": cocktail.glass,
                    "isAlcoholic": cocktail.is_alcoholic
                }

                response = requests.post(
                    self.url,
                    files=files,
                    data=data,
                )
                files["file"][1].close()
                print(response.json())

                if response.ok:
                    logging.info(f"Successfully updated {cocktail.name}.")
                    # Add cocktail name to cache after a successful update
                    self.updated_cocktails_cache.add(cocktail.name)
                else:
                    logging.error(f"Failed to update {cocktail.name}: {response.status_code}")
            except Exception as e:
                logging.error(f"Error updating {cocktail.name}: {e}")
