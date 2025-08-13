from fastapi import APIRouter, Path, Query
from backend.models.models import BlogDetail, BlogListItem
from typing import Annotated, List

blogs = APIRouter(
    prefix="/blogs",
    tags=["blogs"]
)

@blogs.get("/", response_model=list[BlogListItem])
async def get_blogs(
    tags: List[str] = Query(None, description="Filter blogs by tags"),
    limit: int = Query(10, ge=1, le=30),
    offset: int = Query(0, ge=0)):
    """
    Fetch blogs.
    - If `tags` is provided, filter by tags (comma-separated).
    """
    if tags:
        # filter by tags
        pass
    else:
        # return all blogs
        pass

@blogs.get("/{slug}", response_model=BlogDetail)
async def get_blog(slug: Annotated[str, 
                   Path(title="The slug of the blog to get")]):
    pass
