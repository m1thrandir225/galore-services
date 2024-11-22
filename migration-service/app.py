from typing_extensions import Annotated

from migrator import Migrator
from parser import Parser
from fastapi import FastAPI, HTTPException, Depends

from config import Settings
from functools import lru_cache
from fastapi_utilities import repeat_at


@lru_cache
def get_settings():
    return Settings()


app = FastAPI()


origins = [get_settings().api_url]
#
# app.add_middleware(
#     CORSMiddleware,
#     allow_origins=origins,
#     allow_credentials=True,
#     allow_methods=["*"],
#     allow_headers=["*"],
# )

@app.get("/health")
def health_check():
    return {"status": "healthy"}
@app.get("/update-cocktails")
def update_cocktails(settings: Annotated[Settings, Depends(get_settings)]):
    parser = Parser(settings.parser_url, settings.parser_single_url)
    migrator = Migrator(settings.api_url, settings.api_key)
    try:
        cocktails = parser.parse_cocktails()
        migrator.update_cocktails(cocktails)
        return
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
