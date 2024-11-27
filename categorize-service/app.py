from functools import lru_cache
from typing import Annotated

from fastapi import FastAPI, Depends, Request, HTTPException

from categorizer import Categorizer
from config import Settings
from models.detailed_cocktail import DetailedCocktail


@lru_cache
def get_settings():
    return Settings()
app = FastAPI()

@app.get("/health")
def health_check():
    return {"status": "ok"}

@app.post("/categorize")
def categorize_cocktail(
        request: Request,
        cocktail: DetailedCocktail,
        settings: Annotated[Settings, Depends(get_settings)]
):
    request_api = request.headers.get("x-api-key")

    if request_api is None:
        raise HTTPException(status_code=404, detail="API key not found.")
    elif request_api != settings.api_key:
        raise HTTPException(status_code=404, detail="API key not found.")

    categorizer = Categorizer(
        cocktail=cocktail,
        api_service_url=settings.api_service_url,
        api_key=settings.api_key,
    )
    categorizer.categorize_cocktail()


