from pydantic import BaseModel, JsonValue
from datetime import datetime
from uuid import UUID
from typing import List


class TagRead(BaseModel):
    id: int
    name: str


class BlogListItem(BaseModel):
    id: UUID
    title: str
    date: datetime
    slug: str
    description: str
    banner: str
    tags: List[TagRead]


class BlogDetail(BaseModel):
    id: UUID
    title: str
    date: datetime
    modified_at: datetime
    author: str
    slug: str
    description: str
    banner: str
    tags: List[TagRead]
    seo: JsonValue
    openGraph: JsonValue
    twitter: JsonValue
    content_html: str
