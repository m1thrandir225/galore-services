FROM python:3.12

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

# Set the Hugging Face cache directory
ENV HF_HOME=/app/.cache/huggingface

# Preload the Hugging Face model
RUN python -c "from transformers import BertTokenizer, BertModel; \
    BertTokenizer.from_pretrained('bert-base-uncased'); \
    BertModel.from_pretrained('bert-base-uncased')"

COPY . .
COPY .prod.env ./.env


EXPOSE 8000

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]
