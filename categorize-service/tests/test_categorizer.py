import pytest
import responses
from uuid import uuid4
from datetime import datetime

from internal.services.categorizer import CategorizeService
from internal.models import DetailedCocktail
from internal.models.detailed_cocktail import (
    IngredientItem,
    IngredientDTO,
    InstructionsDTO,
)


@pytest.fixture
def mock_cocktail():
    ingredients = IngredientDTO(
        ingredients=[
            IngredientItem(name="Vodka", amount="2 oz"),
            IngredientItem(name="Lime juice", amount="1 oz"),
        ]
    )
    instructions = InstructionsDTO(instructions=["Mix", "Serve"])
    return DetailedCocktail(
        id=uuid4(),
        name="Test Cocktail",
        ingredients=ingredients,
        instructions=instructions,
        image="https://example.com/test.jpg",
        glass="Highball Glass",
        is_alcoholic=True,
        created_at=datetime.now(),
    )


@pytest.fixture
def categorize_service(mock_cocktail):
    return CategorizeService(
        cocktail=mock_cocktail,
        api_service_url="http://categorize-service.local",
        api_key="test-key",
    )


class TestCategorizeByAlcohol:
    @responses.activate
    def test_creates_alcoholic_category_when_not_exists(self, categorize_service):
        responses.add(
            responses.GET,
            "http://categorize-service.local/category/is_alcoholic",
            status=404,
        )
        category_id = str(uuid4())
        responses.add(
            responses.POST,
            "http://categorize-service.local/category",
            json={
                "id": category_id,
                "name": "Alcoholic",
                "tag": "is_alcoholic",
                "created_at": datetime.now().isoformat(),
            },
            status=201,
        )

        responses.add(
            responses.POST,
            "http://categorize-service.local/category_cocktail",
            status=201,
        )

        categorize_service.categorize_by_alcohol()

        assert len(responses.calls) == 3

    @responses.activate
    def test_uses_existing_alcoholic_category(self, categorize_service):
        category_id = str(uuid4())

        responses.add(
            responses.GET,
            "http://categorize-service.local/category/is_alcoholic",
            json={
                "id": category_id,
                "name": "Alcoholic",
                "tag": "is_alcoholic",
                "created_at": datetime.now().isoformat(),
            },
            status=200,
        )

        responses.add(
            responses.POST,
            "http://categorize-service.local/category_cocktail",
            status=201,
        )

        categorize_service.categorize_by_alcohol()

        assert len(responses.calls) == 2


class TestCategorizeByGlass:
    @responses.activate
    def test_creates_glass_category_with_snake_case_tag(self, categorize_service):
        responses.add(
            responses.GET,
            "http://categorize-service.local/category/highball_glass",
            status=404,
        )

        category_id = str(uuid4())
        responses.add(
            responses.POST,
            "http://categorize-service.local/category",
            json={
                "id": category_id,
                "name": "Highball Glass",
                "tag": "highball_glass",
                "created_at": datetime.now().isoformat(),
            },
            status=201,
        )

        responses.add(
            responses.POST,
            "http://categorize-service.local/category_cocktail",
            status=201,
        )

        categorize_service.categorize_by_glass()

        create_call = responses.calls[1]
        assert "highball_glass" in create_call.request.body.decode()  # pyright: ignore[reportAttributeAccessIssue, reportOptionalMemberAccess]


class TestGetFlavoursForCocktail:
    def test_identifies_flavours_from_ingredients(
        self, mock_cocktail, categorize_service
    ):
        flavours = categorize_service._CategorizeService__get_flavours_for_cocktail()
        print(flavours)
        assert "sour" in flavours
        assert "neutral" in flavours
