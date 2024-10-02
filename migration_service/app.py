from warnings import catch_warnings

from parse_service.parser import Parser
from fastapi import FastAPI, HTTPException

app = FastAPI()
parser = Parser()

@app.get("/health")
def health_check():
    return {"status": "healthy"}


@app.get("/update-cocktails")
def update_cocktails():
    try:
        cocktails = parser.parse_cocktails()
        return {"cocktails": cocktails}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
