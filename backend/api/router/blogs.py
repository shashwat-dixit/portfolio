from fastapi import APIRouter

blogs = APIRouter(
    prefix="/blogs",
    tags=["blogs"]
)

@blogs.get("/")
async def get_blogs():
    pass