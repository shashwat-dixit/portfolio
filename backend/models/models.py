from datetime import date
from uuid import uuid4, UUID
from typing import Any

from sqlalchemy import (
    Column, String, ForeignKey, Table
)
from sqlalchemy.dialects.postgresql import UUID as PG_UUID, JSONB
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


# Association table (many-to-many)
blog_post_tag_link = Table(
    "blogposttaglink",
    Base.metadata,
    Column("post_id", PG_UUID(as_uuid=True), ForeignKey("blog.id"), primary_key=True),
    Column("tag_id", PG_UUID(as_uuid=True), ForeignKey("tag.id"), primary_key=True),
)


class Tag(Base):
    __tablename__ = "tag"

    id: Mapped[UUID] = mapped_column(
        PG_UUID(as_uuid=True), primary_key=True, default=uuid4
    )
    name: Mapped[str] = mapped_column(String, unique=True, index=True)

    posts: Mapped[list["Blog"]] = relationship(
        secondary=blog_post_tag_link,
        back_populates="tags"
    )


class Blog(Base):
    __tablename__ = "blog"

    id: Mapped[UUID] = mapped_column(
        PG_UUID(as_uuid=True), primary_key=True, default=uuid4
    )
    title: Mapped[str] = mapped_column(String)
    slug: Mapped[str] = mapped_column(String, index=True, unique=True)
    description: Mapped[str] = mapped_column(String)
    author: Mapped[str] = mapped_column(String, default="Shashwat Dixit")
    published_date: Mapped[date | None]
    last_modified_date: Mapped[date | None]
    image: Mapped[str | None]

    seo: Mapped[dict[str, Any] | None] = mapped_column(JSONB)
    openGraph: Mapped[dict[str, Any] | None] = mapped_column(JSONB)
    twitter: Mapped[dict[str, Any] | None] = mapped_column(JSONB)

    content: Mapped[str] = mapped_column(String)
    file_hash: Mapped[str] = mapped_column(String, index=True)

    tags: Mapped[list[Tag]] = relationship(
        secondary=blog_post_tag_link,
        back_populates="posts"
    )
