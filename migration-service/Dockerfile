FROM python:3.9

WORKDIR /app

COPY requirements.txt .

RUN pip install --upgrade setuptools[core]
RUN pip install --no-cache-dir --upgrade -r requirements.txt

COPY . .
COPY .prod.env ./.env

EXPOSE 8001

CMD ["fastapi", "run", "app.py", "--host", "0.0.0.0", "--port", "8001"]