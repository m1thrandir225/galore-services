import os
from email.header import Header
from functools import lru_cache

from fastapi import FastAPI, HTTPException, Depends
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

@app.get("/health")
def health_check():
    return {"status": "healthy"}

@app.post("/generate-embedding")
async def generate_embedding(data: TextData, x_api_key: str = Header(None), settings: Settings = Depends(get_settings)):
    if x_api_key is None:
        raise HTTPException(status_code=403, detail="X-Api-Key not provided")
    elif x_api_key != settings.api_key:
        raise HTTPException(status_code=403, detail="X-Api-Key not provided")

    try:
        embedding = model.get_embeddings(data.text)
        return {"embedding": embedding.tolist()}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
