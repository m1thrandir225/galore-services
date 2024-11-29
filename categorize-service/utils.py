import json
from datetime import datetime
# Text to snake case for typically used for tags
def to_snake_case(text: str) -> str:
    return text.lower().replace(" ", "_").replace("-", "_")

def parse_time(time: str) -> datetime:
    return datetime.fromisoformat(time)

def is_json(text: str) -> bool:
    try:
        json.loads(text)
    except ValueError as e:
        return False
    return True