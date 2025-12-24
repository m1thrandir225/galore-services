FROM python:3.14-slim as builder

WORKDIR /app

RUN pip install uv 

COPY ./pyproject.toml ./uv.lock ./ 

RUN uv sync 
RUN uv cache clean

FROM python:3.14-slim as runner 

WORKDIR /app

COPY --from=builder /app/.venv /app/.venv
COPY . .

ENV PATH="/app/.venv/bin:$PATH"

RUN ["uvicorn", "main:app", "--proxy-headers", "--port", "8003"]
