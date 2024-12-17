FROM python:3.9

WORKDIR /app

COPY requirements.txt .

RUN pip install --upgrade setuptools[core]
RUN pip install --no-cache-dir --upgrade -r requirements.txt

COPY . .
EXPOSE 8003

CMD ["fastapi", "run", "app.py", "--proxy-headers","--port", "8003"]