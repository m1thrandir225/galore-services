import json
import time
from dataclasses import asdict

import logging
from typing import List
import uuid

import requests
from requests import HTTPError, request

import utils
from models.detailed_cocktail import DetailedCocktail


class Migrator:
    def __init__(self, url: str, api_key: str, max_retries: int = 5):
        self.url = url
        self.api_key = api_key
        self.updated_cocktails_cache = (
            set()
        )  # Cache to store successfully updated cocktail names or IDs
        self.max_retries = max_retries

    async def wait_for_service(self):
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
            except HTTPError as e:
                logging.warning(f"Service not ready, attempt {attempt + 1}/{self.max_retries}")
                time.sleep(5)  # Wait 5 seconds between retries

        raise Exception("Could not connect to the service after multiple attempts")


    async def update_cocktails(self, cocktails: List[DetailedCocktail]):
        """Updates cocktails by sending data to the given URL with a caching mechanism."""
        await self.wait_for_service()

        for cocktail in cocktails:
            # Use cocktail name or ID as a unique identifier for the cache
            if cocktail.name in self.updated_cocktails_cache:
                logging.info(f"Skipping {cocktail.name}, already updated.")
                return

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

                response = requests.post(
                    self.url + "/migration/cocktails",
                    files=files,
                    data=data,
                    headers={"x-api-key": self.api_key},
                )
                files["file"][1].close()

                logging.info(f"Successfully updated {cocktail.name}.")

                # Add cocktail name to cache after a successful update
                self.updated_cocktails_cache.add(cocktail.name)

            except Exception as e:
                logging.error(f"Error updating {cocktail.name}: {e}")