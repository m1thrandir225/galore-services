import pytest
from fastapi.testclient import TestClient
from uuid import uuid4
from datetime import datetime

from main import app, get_settings
from internal.config import Settings


class TestSettings(Settings):
    api_service_url: str = "http://categorize-service.local"
    api_key: str = "test-api-key"


def get_test_settings():
    return TestSettings()


@pytest.fixture
def client():
    app.dependency_overrides[get_settings] = get_test_settings
    with TestClient(app) as c:
        yield c
    app.dependency_overrides.clear()


@pytest.fixture
def sample_cocktail_data():
    return {
        "id": str(uuid4()),
        "name": "Margarita",
        "ingredients": {
            "ingredients": [
                {"name": "Tequila", "amount": "2 oz"},
                {"name": "Lime juice", "amount": "1 oz"},
                {"name": "Simple syrup", "amount": "0.5 oz"},
            ]
        },
        "instructions": {
            "instructions": ["Shake all ingredients with ice", "Strain into glass"]
        },
        "image": "https://example.com/margarita.jpg",
        "glass": "Cocktail Glass",
        "is_alcoholic": True,
        "created_at": datetime.now().isoformat(),
    }


@pytest.fixture
def api_headers():
    return {"x-api-key": "test-api-key"}
