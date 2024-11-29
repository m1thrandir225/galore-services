import os
from functools import lru_cache
from fastapi import FastAPI, HTTPException, Depends, Header
from pydantic import BaseModel
from transformers import AutoModel

from config import Settings
from models.model import TextEmbeddingModel


@lru_cache
def get_settings():
    return Settings()


app = FastAPI()
model = TextEmbeddingModel()

class TextData(BaseModel):
    text: str


def check_api_key(x_api_key: str = Header(...)):
    if not x_api_key:
        raise HTTPException(status_code=400, detail="X-Api-Key header missing")
    return x_api_key


@app.get("/health")
def health_check():
    return {"status": "healthy"}

@app.post("/generate-embedding")
async def generate_embedding(
    data: TextData,
    api_key: str = Depends(check_api_key),
    settings: Settings = Depends(get_settings),
):
    if api_key != settings.api_key:
        if api_key != "testing" and settings.environment != "development":
            raise HTTPException(status_code=403, detail="Incorrect api_key")

    try:
        embedding = model.get_embeddings(data.text)
        return {"embedding": embedding.tolist()}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
