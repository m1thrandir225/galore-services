import json
import logging
import time
import uuid
from dataclasses import asdict

import requests
from requests import HTTPError

from ..models import DetailedCocktail


class MigrateService:
    def __init__(self, url: str, api_key: str, max_retries: int = 5):
        self.url: str = url
        self.api_key: str = api_key
        self.max_retries: int = max_retries

    async def wait_for_service(self) -> bool:
        """
        Wait for the service to be available
        """
        for attempt in range(self.max_retries):
            try:
                # Use a health check endpoint if available
                response = requests.get(f"{self.url}/health", timeout=5)
                response.raise_for_status()
                logging.info("Service is ready!")
                return True
            except HTTPError:
                logging.warning(
                    f"Service not ready, attempt {attempt + 1}/{self.max_retries}"
                )
                time.sleep(5)  # Wait 5 seconds between retries

        raise Exception("Could not connect to the service after multiple attempts")

    async def update_cocktails(self, cocktails: list[DetailedCocktail]) -> None:
        """Updates cocktails by sending data to the given URL with a caching mechanism."""
        _ = await self.wait_for_service()

        for cocktail in cocktails:
            try:
                filename = f"{str(uuid.uuid4())}.jpg"
                files = {
                    "file": (filename, open(cocktail.image, "rb"), "image/jpg"),
                }

                data = {
                    "name": cocktail.name,
                    "ingredients": json.dumps(asdict(cocktail.ingredients)),
                    "instructions": cocktail.instructions,
                    "glass": cocktail.glass,
                    "isAlcoholic": cocktail.isAlcoholic,
                }

                _ = requests.post(
                    self.url + "/migration/cocktails",
                    files=files,
                    data=data,
                    headers={"x-api-key": self.api_key},
                )
                files["file"][1].close()

                logging.info(f"Successfully updated {cocktail.name}.")

            except Exception as e:
                logging.error(f"Error updating {cocktail.name}: {e}")

