import json
import time
from os.path import isfile
from typing import Dict, List, Optional
import requests

import utils
from models.detailed_cocktail import DetailedCocktail
import logging

# Set up logging for the cron job
logging.basicConfig(level=logging.INFO)

class Parser:
    def __init__(self, url: str, single_cocktail_url: str):
        self.url: str = url
        self.single_cocktail_url = single_cocktail_url
        self.cocktail_list_cache = None  # Cache for cocktail list
        self.cocktail_details_cache = {}  # Cache for individual cocktail details

    def get_cocktails(self) -> Optional[List]:
        """Fetch the list of cocktails, with error handling and caching."""
        if self.cocktail_list_cache is not None:
            logging.info("Returning cached cocktail list.")
            return self.cocktail_list_cache

        try:
            response = requests.get(self.url, timeout=10)
            response.raise_for_status()  # Raises an error for bad status codes (4xx, 5xx)
            json_data: List = response.json()["drinks"]
            self.cocktail_list_cache = json_data  # Cache the result
            logging.info("Successfully fetched cocktail list.")
            return json_data
        except requests.exceptions.RequestException as e:
            logging.error(f"Failed to fetch cocktail list: {e}")
            return None


    async def parse_cocktails(self) -> List["DetailedCocktail"]:
        """Fetch detailed data for each cocktail, using caching."""
        cocktails: List[DetailedCocktail] = []

        data = self.get_cocktails()
        if data is None:
            return cocktails  # Return empty list if there was an error fetching data

        for item in data:

            drink_id = item["idDrink"]

            if drink_id in self.cocktail_details_cache:
                detailed_cocktail = self.cocktail_details_cache[drink_id]
                logging.info(f"Using cached data for drink ID {drink_id}.")
            else:
                try:
                    headers = {"Accept": "application/json"}
                    response = requests.get(
                        f"{self.single_cocktail_url}{drink_id}",
                        headers=headers,
                        timeout=10,
                    )
                    response.raise_for_status()
                    single_cocktail_data = response.json()["drinks"][0]
                    detailed_cocktail = DetailedCocktail(
                        id=single_cocktail_data["idDrink"],
                        name=single_cocktail_data["strDrink"],
                        ingredients=utils.format_ingredients(single_cocktail_data),
                        instructions=single_cocktail_data["strInstructions"],
                        image=utils.download_image(
                            single_cocktail_data["strDrinkThumb"]
                        ),
                        glass=single_cocktail_data["strGlass"],
                        isAlcoholic=utils.has_alc(single_cocktail_data["strAlcoholic"]),
                    )

                    # Cache the cocktail details
                    self.cocktail_details_cache[drink_id] = detailed_cocktail
                    logging.info(f"Fetched details for drink ID {drink_id}.")
                except requests.exceptions.RequestException as e:
                    logging.error(
                        f"Failed to fetch details for drink ID {drink_id}: {e}"
                    )
                    return cocktails

                cocktails.append(detailed_cocktail)
                time.sleep(0.5)  # Respect API rate limits

        logging.info(f"Fetched details for {len(cocktails)} cocktails.")
        return cocktails
