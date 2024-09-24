import os

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from transformers import AutoModel

from models.model import TextEmbeddingModel

app = FastAPI()
model = TextEmbeddingModel()

class TextData(BaseModel):
    text: str

@app.get("/health")
def health_check():
    return {"status": "healthy"}

@app.post("/generate-embedding")
async def generate_embedding(data: TextData):
    try:
        embedding = model.get_embeddings(data.text)
        return {"embedding": embedding.tolist()}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))