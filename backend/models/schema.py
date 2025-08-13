from typing import Optional, List
from datetime import datetime
from uuid import UUID
from sqlmodel import SQLModel, Field, Relationship
from pydantic import JsonValue


class BlogTagLink(SQLModel, table=True):
    blog_id: UUID = Field(foreign_key="blog.id", primary_key=True)
    tag_id: int = Field(foreign_key="tag.id", primary_key=True)


class Tag(SQLModel, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)
    name: str

    blogs: List["Blog"] = Relationship(back_populates="tags", link_model=BlogTagLink)


class Blog(SQLModel, table=True):
    id: UUID = Field(primary_key=True)
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

    tags: List[Tag] = Relationship(back_populates="blogs", link_model=BlogTagLink)
    content: Optional["BlogContent"] = Relationship(back_populates="blog")


class BlogContent(SQLModel, table=True):
    blog_id: UUID = Field(foreign_key="blog.id", primary_key=True)
    content_raw: str  # original markdown from Obsidian
    content_html: str  # parsed HTML with CDN images

    blog: Optional[Blog] = Relationship(back_populates="content")
