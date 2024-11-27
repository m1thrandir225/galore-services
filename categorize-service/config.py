from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    app_name: str = "Galore-Categorize-Service"
    api_service_url: str
    api_key: str
    model_config = SettingsConfigDict(env_file=".env")