from typing import Optional, List
from pydantic import JsonValue
from datetime import datetime
from uuid import UUID
from sqlmodel import SQLModel, Field, Relationship

class BlogTagLink(SQLModel, table=True):
    blog_id: UUID = Field(foreign_key="blog.id", primary_key=True)
    tag_id: int = Field(foreign_key="tag.id", primary_key=True)


class Tag(SQLModel, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)
    name: str

    blogs: List["Blog"] = Relationship(
        back_populates="tags", link_model=BlogTagLink
    )


class Blog(SQLModel, table=True):
    id: UUID
    title: str
    date: datetime
    modified_at: datetime
    author: str
    slug: str
    description: str
    banner: str
    seo: JsonValue
    openGraph: JsonValue
    twitter: JsonValue

    tags: List[Tag] = Relationship(
        back_populates="blogs", link_model=BlogTagLink
    )


class BlogContent(SQLModel, table=True):
    blog_id: UUID = Field(foreign_key="blog.id", primary_key=True)
    content: str  # store Markdown directly as string
    blog: Optional["Blog"] = Relationship(back_populates="content")