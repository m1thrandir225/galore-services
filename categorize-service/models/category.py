from datetime import datetime
import uuid

from pydantic import BaseModel

class CategoryDTO(BaseModel):
    name: str
    tag: str

class Category(BaseModel):
    id: uuid.UUID
    name: str
    tag: str
    created_at: datetime