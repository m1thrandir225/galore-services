import os
from functools import lru_cache
from fastapi import FastAPI, HTTPException, Depends, Header
from pydantic import BaseModel
from starlette import status
from starlette.requests import Request
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

@app.post("/generate-embedding")
async def generate_embedding(
    data: TextData,
    api_key: str = Depends(validate_api_key),
    settings: Settings = Depends(get_settings)
):
    try:
        embedding = model.get_embeddings(data.text)
        return {"embedding": embedding.tolist()}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
