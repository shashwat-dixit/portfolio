from datetime import date
from uuid import UUID, uuid4
from typing import Any
from sqlmodel import SQLModel, Field, Column, Relationship
from sqlalchemy.dialects.postgresql import JSONB


class BlogPostTagLink(SQLModel, table=True):
    post_id: UUID = Field(foreign_key="blog.id", primary_key=True)
    tag_id: UUID = Field(foreign_key="tag.id", primary_key=True)


class Tag(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True)
    name: str = Field(unique=True, index=True)

    posts: list["Blog"] = Relationship(back_populates="tags", link_model=BlogPostTagLink)


class Blog(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True)
    title: str
    slug: str = Field(index=True, unique=True)
    description: str
    author: str = Field(default="Shashwat Dixit")
    published_date: date | None = Field(default=None)
    last_modified_date: date | None = Field(default=None)
    image: str | None = None
    seo: dict[str, Any] | None = Field(default=None, sa_column=Column(JSONB))
    openGraph: dict[str, Any] | None = Field(default=None, sa_column=Column(JSONB))
    twitter: dict[str, Any] | None = Field(default=None, sa_column=Column(JSONB))
    content: str
    tags: list[Tag] = Relationship(back_populates="posts", link_model=BlogPostTagLink)
    file_hash: str = Field(nullable=False, index=True)
    # Not doing this now, if needed will add it if perf is not satisfactory.
    # Stores Pre-rendered HTML
    # content_html: str | None = None
