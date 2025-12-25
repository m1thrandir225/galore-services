FROM python:3.14-slim

WORKDIR /app

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

COPY ./pyproject.toml ./uv.lock ./

RUN uv pip install --system -r pyproject.toml

# Set the Hugging Face cache directory
ENV HF_HOME=/app/.cache/huggingface

# Preload the Hugging Face model
RUN python -c "from transformers import BertTokenizer, BertModel; \
  BertTokenizer.from_pretrained('bert-base-uncased'); \
  BertModel.from_pretrained('bert-base-uncased')"

COPY . .

EXPOSE 8000

CMD ["fastapi", "run", "--host", "0.0.0.0" , "--proxy-headers","--port", "8000"]
