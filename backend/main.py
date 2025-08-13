from fastapi import FastAPI
from backend.api.router.blogs import blogs

app = FastAPI(title="Portfolio Blog API")
app.include_router(blogs, prefix="/api")

@app.get("/")
async def root():
    return {"message": "Welcome to Portfolio Blog API"}