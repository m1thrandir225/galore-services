import logging
from typing import Set

import requests
from requests import HTTPError

from models.category import Category, CategoryDTO
from models.detailed_cocktail import DetailedCocktail
from utils import to_snake_case
from models.flavour_map import flavor_map

"""
Categorizer
"""
class Categorizer:
    def __init__(self, cocktail: DetailedCocktail, api_service_url: str, api_key: str):
        self.cocktail = cocktail
        self.api_service_url = api_service_url
        self.api_key =  api_key
    """
        Should categorize the cocktail based on 3 categories: 
        1. Based on is it alcoholic
        2. Based on glass type
        3. Based on flavour
        
        For each of those it should check if the category already exists, if it does create a
        category_cocktail object with the current cocktail. If it does not then it should create the category
        and them create the category_cocktail object.
    """

    def __category_exists(self, category_tag: str) -> Category or False or None:
        """
        Checks to see if the given category exists
        1. If the category exists, return it.
        2. If the category does not exist, return False.
        3. If there is an error, return None
        """

        try:
            response = requests.get(
                self.api_service_url + "/category/" + category_tag,
                headers={"x-api-key":self.api_key},
            )
            response.raise_for_status()

            j = response.json()
            return Category(**j)
        except HTTPError as e:
            code = e.response.status_code
            if code != 404:
                logging.error(f"Exception while checking if category exists: {e}")
                return None
            return False

    def __create_category(self, category: CategoryDTO) -> Category:
        """Creates a new category"""
        try:
            response = requests.post(
                self.api_service_url + "/category",
                json={
                    "name": category.name,
                    "tag": category.tag,
                },
                headers={"x-api-key":self.api_key},
            )
            if not response.ok:
                raise Exception(f"Error status code: {response.status_code}")

            logging.info(f"Created new category {category.name}")
            j = response.json()
            return Category(**j)
        except HTTPError as e:
            logging.error(f"Failed to create new category {category.name}, error: {e.response.text}")

    def __create_category_cocktail(self, category_id: str):
        """Creates a category_cocktail object with an existing category and cocktail"""
        try:
            data = {
                "category_id": category_id,
                "cocktail_id": str(self.cocktail.id)
            }
            response = requests.post(
                self.api_service_url + "/category_cocktail",
                json=data,
                headers={"x-api-key":self.api_key},
            )
            response.raise_for_status()

            logging.info(f"Created new cocktail_category object with category_id:{category_id}, cocktail_id:{self.cocktail.id}")
        except HTTPError as e:
            logging.info(f"Failed to create ne cocktail_category object with category_id:{category_id}, error: {e.response.text}")

    def categorize_by_glass(self):
        """Categorizes the cocktail by glass"""
        glass_type =  self.cocktail.glass
        tag_glass_type = to_snake_case(glass_type.lower())

        category = self.__category_exists(category_tag=tag_glass_type)
        if isinstance(category, Category):
            """
            Category exists
            1. Create category_cocktail object
            """
            self.__create_category_cocktail(category_id=str(category.id))
        elif category is False:
            """ 
            Category does not exist
            1. Create the category
            2. Create category_cocktail object
            """
            category_dto = CategoryDTO(
                name=glass_type,
                tag=tag_glass_type,
            )
            new_category = self.__create_category(category=category_dto)
            self.__create_category_cocktail(category_id=str(new_category.id))

    def __get_flavours_for_cocktail(self) -> Set[str]:
        flavours = set()
        for ingredient in self.cocktail.ingredients.ingredients:
            name = ingredient.name.lower()
            for key, flavour in flavor_map.items():
                if flavour in name:
                    flavours.add(flavour)

        return flavours

    def categorize_by_flavour(self):
        """Categorizes the cocktail by flavour"""
        cocktail_flavours = self.__get_flavours_for_cocktail()

        for flavour in cocktail_flavours:
            tag = to_snake_case(flavour.lower())
            name = flavour.capitalize()

            category = self.__category_exists(category_tag=tag)
            if isinstance(category, Category):
                self.__create_category_cocktail(category_id=str(category.id))
            elif category is False:
                category_dto = CategoryDTO(
                    name=name,
                    tag=tag,
                )
                new_category = self.__create_category(category=category_dto)
                self.__create_category_cocktail(category_id=str(new_category.id))

    def categorize_by_alcohol(self):
        """Categorizes the cocktail by alcohol"""
        is_alcoholic = self.cocktail.is_alcoholic
        alcoholic_name = "Alcoholic"
        alcoholic_tag = "is_alcoholic"
        not_alcoholic_name = "Not Alcoholic"
        not_alcoholic_tag = "is_not_alcoholic"
        if is_alcoholic:
            category = self.__category_exists(category_tag=alcoholic_tag)
            if isinstance(category, Category):
                self.__create_category_cocktail(category_id=str(category.id))
            elif category is False:
                category_dto = CategoryDTO(
                    name=alcoholic_name,
                    tag=alcoholic_tag,
                )
                new_category = self.__create_category(category=category_dto)
                self.__create_category_cocktail(category_id=str(new_category.id))
        else:
            category = self.__category_exists(category_tag=not_alcoholic_tag)

            if isinstance(category, Category):
                self.__create_category_cocktail(category_id=str(category.id))
            elif category is False:
                category_dto = CategoryDTO(
                    name=not_alcoholic_name,
                    tag=not_alcoholic_tag,
                )
                new_category = self.__create_category(category=category_dto)
                self.__create_category_cocktail(category_id=str(new_category.id))

