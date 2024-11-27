from dis import Instruction
from typing import List

from pydantic import BaseModel

class InstructionsDTO(BaseModel):
    instructions: List[str]