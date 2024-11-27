from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    app_name: str = "Galore-Migration-Service"
    api_url: str
    api_key: str
    parser_url: str
    parser_single_url: str
    model_config = SettingsConfigDict(env_file=".env")
