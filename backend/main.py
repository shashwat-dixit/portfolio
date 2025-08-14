from fastapi import FastAPI
from api.router.blogs import blogs
from core.database import create_db_and_tables


app = FastAPI(title="Portfolio Blog API")
app.include_router(blogs, prefix="/api")

create_db_and_tables()

@app.get("/")
async def root():
    return {"message": "Welcome to Portfolio Blog API"}