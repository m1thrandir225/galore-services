from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    app_name: str = "Galore-Embedding-Service"
    api_key: str
    model_config = SettingsConfigDict(env_file=".env")

