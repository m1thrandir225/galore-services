FROM python:3.9

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

EXPOSE 8000

CMD ["fastapi", "run", "app.py", "--proxy-headers","--port", "8000"]
