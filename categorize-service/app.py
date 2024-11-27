import json
import logging
import uuid
from functools import lru_cache
from typing import Annotated

from fastapi import FastAPI, Depends, Request, HTTPException, Body, status

from categorizer import Categorizer
from config import Settings
from models.detailed_cocktail import DetailedCocktail
from models.ingredient import Ingredient
from models.ingredient_dto import IngredientDTO
from models.instructions_dto import InstructionsDTO


@lru_cache
def get_settings():
    return Settings()
app = FastAPI()

@app.get("/health")
def health_check():
    return {"status": "ok"}

async def validate_api_key(request: Request, settings: Settings = Depends(get_settings)):
    x_api_key = request.headers.get("x-api-key")
    if not x_api_key:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Missing x-api-key in headers"
        )
    if x_api_key != settings.api_key:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Invalid API key"
        )
    return x_api_key

@app.post("/categorize")
async def categorize_cocktail(
        request: Request,
        api_key: str = Depends(validate_api_key),
        settings: Settings = Depends(get_settings)
):
    r_json = await request.json()
    print(settings.api_service_url)
    ingredients_json = (r_json["ingredients"]["ingredients"])
    print(ingredients_json)
    ingredients = []
    for ingredient in ingredients_json:
        ingredients.append(Ingredient(
            name=ingredient["name"],
            amount=ingredient["amount"],
        ))
    ingredient_dto = IngredientDTO(
        ingredients=ingredients,
    )
    instructions_dto = InstructionsDTO(
        instructions= json.loads(r_json["instructions"])["instructions"]
    )

    cocktail = DetailedCocktail(
        name=r_json["name"],
        ingredients=ingredient_dto,
        instructions=instructions_dto,
        embedding=r_json["embedding"],
        is_alcoholic=r_json["is_alcoholic"],
        image=r_json["image"],
        created_at=r_json["created_at"],
        glass=r_json["glass"],
        id=r_json["id"],
    )
    categorizer = Categorizer(
        cocktail=cocktail,
        api_service_url=settings.api_service_url,
        api_key=settings.api_key,
    )
    await categorizer.categorize_cocktail()


