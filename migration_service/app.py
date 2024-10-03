from typing_extensions import Annotated

from migrator import Migrator
from parser import Parser
from fastapi import FastAPI, HTTPException, Depends

from config import Settings
from functools import lru_cache
@lru_cache
def get_settings():
    return Settings()

app = FastAPI()
parser = Parser()

origins = [
    get_settings().api_url
]
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
    migrator = Migrator(settings.api_url)
    try:
        cocktails = parser.parse_cocktails()
        migrator.update_cocktails(cocktails)
        return
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
