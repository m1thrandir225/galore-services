from functools import lru_cache
from typing import Annotated

from fastapi import Depends, FastAPI, HTTPException, Request, status
from .internal.config import Settings
from .internal.models import DetailedCocktail
from .internal.services import CategorizeService


@lru_cache
def get_settings() -> Settings:
    return Settings()  # pyright: ignore[reportCallIssue]


app = FastAPI()


@app.get("/health")
def health_check():
    return {"status": "ok"}


async def validate_api_key(
    request: Request, settings: Annotated[Settings, Depends(get_settings)]
):
    x_api_key = request.headers.get("x-api-key")
    if not x_api_key:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Missing x-api-key in headers",
        )
    if x_api_key != settings.api_key:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN, detail="Invalid API key"
        )
    return x_api_key


@app.post("/categorize")
async def categorize_cocktail(
    cocktail: DetailedCocktail,
    api_key: Annotated[str, Depends(validate_api_key)],
    settings: Annotated[Settings, Depends(get_settings)],
):
    categorizer = CategorizeService(
        cocktail=cocktail,
        api_service_url=settings.api_service_url,
        api_key=settings.api_key,
    )

    categorizer.categorize_by_alcohol()
    categorizer.categorize_by_flavour()
    categorizer.categorize_by_glass()
