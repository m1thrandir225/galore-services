from starlette import status
from starlette.requests import Request
from typing_extensions import Annotated

from migrator import Migrator
from parser import Parser
from fastapi import FastAPI, HTTPException, Depends, Header, BackgroundTasks

from config import Settings
from functools import lru_cache


@lru_cache
def get_settings():
    return Settings()

app = FastAPI()

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

@app.get("/health")
def health_check():
    return {"status": "healthy"}


@app.get("/update-cocktails", status_code=status.HTTP_202_ACCEPTED)
async def update_cocktails(
    settings: Annotated[Settings, Depends(get_settings)],
    background_tasks: BackgroundTasks,
    # api_key: str = Depends(validate_api_key),
):
    parser = Parser(settings.parser_url, settings.parser_single_url)
    migrator = Migrator(url=settings.api_url, api_key=settings.api_key)
    cocktails = await parser.parse_cocktails()
    background_tasks.add_task(migrator.update_cocktails, cocktails)
    return {"status": "Migration started in the background."}
