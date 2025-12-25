import pytest


class TestHealthCheck:
    def test_health_check_returns_ok(self, client):
        response = client.get("/health")
        assert response.status_code == 200
        assert response.json() == {"status": "ok"}


class TestCategorizeEndpoint:
    def test_missing_api_key_returns_401(self, client, sample_cocktail_data):
        response = client.post("/categorize", json=sample_cocktail_data)
        assert response.status_code == 401
        assert "Missing x-api-key" in response.json()["detail"]

    def test_invalid_api_key_returns_403(self, client, sample_cocktail_data):
        response = client.post(
            "/categorize",
            json=sample_cocktail_data,
            headers={"x-api-key": "wrong-key"},
        )
        assert response.status_code == 403
        assert "Invalid API key" in response.json()["detail"]

    def test_invalid_cocktail_data_returns_422(self, client, api_headers):
        response = client.post(
            "/categorize",
            json={"invalid": "data"},
            headers=api_headers,
        )
        assert response.status_code == 403
