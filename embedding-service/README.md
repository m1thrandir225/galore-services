# Embedding Service

The embedding-service is meant for generating embeddings to enable
easier vector search for the cocktails.

It's built using:
- Python: Language 
- FastAPI: HTTP Server
- BertTokenizer: The tokenizer that generates the embeddings

# Run Locally 
This app is not meant to be run as a standalone service. It works in collaboration with
the `api-service`.

## Prerequisites
- Set up the virtual environment using the provided venv
- Create a `.env` file from the provided `.env.example` file.

## Steps
1. Run the app using the cli :
   1. For production: `fastapi run`
   2. For development: `fastapi dev`