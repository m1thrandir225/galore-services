import json
import time
from os.path import isfile
from typing import Dict, List, Optional
import requests
from models.detailed_cocktail import DetailedCocktail
import logging

# Set up logging for the cron job
logging.basicConfig(level=logging.INFO)

class Parser:
    def __init__(self):
        self.url: str = (
            "https://www.thecocktaildb.com/api/json/v1/1/filter.php?c=Cocktail"
        )
        self.single_cocktail_url = (
            "https://www.thecocktaildb.com/api/json/v1/1/lookup.php?i="
        )
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

    def parse_cocktails(self) -> List["DetailedCocktail"]:
        """Fetch detailed data for each cocktail, using caching."""
        cocktails: List["DetailedCocktail"] = []

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
                    detailed_cocktail = DetailedCocktail(single_cocktail_data)

                    # Cache the cocktail details
                    self.cocktail_details_cache[drink_id] = detailed_cocktail
                    logging.info(f"Fetched details for drink ID {drink_id}.")
                except requests.exceptions.RequestException as e:
                    logging.error(
                        f"Failed to fetch details for drink ID {drink_id}: {e}"
                    )
                    continue  # Skip to the next drink on failure

            cocktails.append(detailed_cocktail)
            time.sleep(0.5)  # Respect API rate limits
            break

        logging.info(f"Fetched details for {len(cocktails)} cocktails.")
        return cocktails
