import json
from datetime import datetime


def to_snake_case(text: str) -> str:
    """
    Text to snake case typically used for tags
    """
    return text.lower().replace(" ", "_").replace("-", "_")


def parse_time_iso(time: str) -> datetime:
    """
    Parses an iso string to a datetime object
    """
    return datetime.fromisoformat(time)


def is_json(text: str) -> bool:
    """
    Checks if a given string is a valid JSON
    """
    try:
        json.loads(text)
    except ValueError as e:
        return False
    return True
