from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    app_name: str = "Galore-Migration-Service"
    api_url: str

    model_config = SettingsConfigDict(env_file=".env")