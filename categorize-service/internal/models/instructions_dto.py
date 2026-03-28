from pydantic import BaseModel


class InstructionsDTO(BaseModel):
    instructions: list[str]

