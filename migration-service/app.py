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
async def update_cocktails(settings: Annotated[Settings, Depends(get_settings)]):
    print(settings.api_key, settings.api_url)

    parser = Parser(settings.parser_url, settings.parser_single_url)
    migrator = Migrator(url=settings.api_url, api_key=settings.api_key)
    try:
        cocktails = await parser.parse_cocktails()
        await migrator.update_cocktails(cocktails)
        return {"status": "success"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
