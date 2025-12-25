import pytest
from datetime import datetime

from internal.config.utils import to_snake_case, parse_time_iso, is_json


class TestToSnakeCase:
    def test_simple_text(self):
        assert to_snake_case("Hello World") == "hello_world"

    def test_with_hyphens(self):
        assert to_snake_case("Cocktail-Glass") == "cocktail_glass"

    def test_already_lowercase(self):
        assert to_snake_case("already lowercase") == "already_lowercase"

    def test_mixed_separators(self):
        assert to_snake_case("Some-Text With Spaces") == "some_text_with_spaces"


class TestParseTimeIso:
    def test_valid_iso_string(self):
        result = parse_time_iso("2024-12-24T10:30:00")
        assert isinstance(result, datetime)
        assert result.year == 2024
        assert result.month == 12
        assert result.day == 24

    def test_invalid_iso_string(self):
        with pytest.raises(ValueError):
            parse_time_iso("not-a-date")


class TestIsJson:
    def test_valid_json_object(self):
        assert is_json('{"key": "value"}') is True

    def test_valid_json_array(self):
        assert is_json("[1, 2, 3]") is True

    def test_invalid_json(self):
        assert is_json("not json") is False

    def test_empty_string(self):
        assert is_json("") is False
