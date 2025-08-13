from fastapi import FastAPI
import api.blogs

app = FastAPI(title="Portfolio Blog API")
app.include_router(api.blogs.blogs, prefix="/api")

@app.get("/")
async def root():
    return {"message": "Welcome to Portfolio Blog API"}