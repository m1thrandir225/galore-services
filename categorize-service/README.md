# Categorize Service

The categorize service is meant for categorizing cocktails based upon three
of their properties: are they alcoholic, flavour and  by glass type. 

It's built using:
- Python: Language
- FastAPI: HTTP Server

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