FROM python:3.14-slim

WORKDIR /app

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

COPY ./pyproject.toml ./uv.lock ./ 

RUN uv pip install --system -r pyproject.toml 

COPY . . 

EXPOSE 8001

CMD ["fastapi", "run", "--host", "0.0.0.0", "--port", "8001", "--proxy-headers"]
