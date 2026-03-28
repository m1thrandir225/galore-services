import time
import requests
from ..config import format_ingredients, download_image, has_alc
from ..models import DetailedCocktail
import logging

logging.basicConfig(level=logging.INFO)


class ParseService:
    """
    Parser
    """

    def __init__(self, url: str, single_cocktail_url: str):
        self.url: str = url
        self.single_cocktail_url: str = single_cocktail_url
        self.cocktail_details_cache: dict[
            str, DetailedCocktail
        ] = {}  # Cache for individual cocktail details

    def get_cocktails(self) -> list[dict[str, str]] | None:
        try:
            response = requests.get(self.url, timeout=10)
            response.raise_for_status()  # Raises an error for bad status codes (4xx, 5xx)
            json_data: list[dict[str, str]] = response.json()["drinks"]
            logging.info("Successfully fetched cocktail list.")
            return json_data
        except requests.exceptions.RequestException as e:
            logging.error(f"Failed to fetch cocktail list: {e}")
            return None

    async def parse_cocktails(self) -> list[DetailedCocktail]:
        """Fetch detailed data for each cocktail, using caching."""
        cocktails: list[DetailedCocktail] = []

        data = self.get_cocktails()
        if data is None:
            return cocktails  # Return empty list if there was an error fetching data

        for item in data:
            drink_id = item["idDrink"]
            try:
                headers = {"Accept": "application/json"}
                response = requests.get(
                    f"{self.single_cocktail_url}{drink_id}",
                    headers=headers,
                    timeout=10,
                )
                response.raise_for_status()
                single_cocktail_data: dict[str, str] = response.json()["drinks"][0]

                detailed_cocktail = DetailedCocktail(
                    id=single_cocktail_data["idDrink"],
                    name=single_cocktail_data["strDrink"],
                    ingredients=format_ingredients(single_cocktail_data),
                    instructions=single_cocktail_data["strInstructions"],
                    image=download_image(single_cocktail_data["strDrinkThumb"]),
                    glass=single_cocktail_data["strGlass"],
                    isAlcoholic=has_alc(single_cocktail_data["strAlcoholic"]),
                )

                # Cache the cocktail details
                self.cocktail_details_cache[drink_id] = detailed_cocktail
                logging.info(f"Fetched details for drink ID {drink_id}.")
                cocktails.append(detailed_cocktail)
                time.sleep(0.5)  # Respect API rate limits
            except requests.exceptions.RequestException as e:
                logging.error(f"Failed to fetch details for drink ID {drink_id}: {e}")
                return cocktails

        logging.info(f"Fetched details for {len(cocktails)} cocktails.")
        return cocktails
