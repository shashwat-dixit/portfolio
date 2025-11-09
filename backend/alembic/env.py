import asyncio
from logging.config import fileConfig
from sqlalchemy.ext.asyncio import create_async_engine
from alembic import context
import sqlmodel
from sqlmodel import SQLModel
from dotenv import load_dotenv
import os

load_dotenv()

# ✅ import your models (adjust path if needed)
from models.models import Blog, Tag, BlogPostTagLink  # <-- your models

# Alembic Config
config = context.config

if config.config_file_name is not None:
    fileConfig(config.config_file_name)

target_metadata = SQLModel.metadata

DATABASE_URL = os.getenv("DATABASE_URL")
if not DATABASE_URL:
    raise RuntimeError("DATABASE_URL must be set")

config.set_main_option("sqlalchemy.url", DATABASE_URL)


def run_migrations_offline():
    """Offline mode — generates SQL but does not connect."""
    context.configure(
        url=DATABASE_URL,
        target_metadata=target_metadata,
        literal_binds=True,
        compare_type=True,
    )

    with context.begin_transaction():
        context.run_migrations()

async def run_migrations_online():
    connectable = create_async_engine(DATABASE_URL, future=True)

    async with connectable.connect() as connection:
        await connection.run_sync(
            lambda sync_conn: context.configure(
                connection=sync_conn,
                target_metadata=target_metadata,
                compare_type=True,
            )
        )

        await connection.run_sync(lambda sync_conn: context.run_migrations())


def run_migrations_online_sync():
    asyncio.run(run_migrations_online())


if context.is_offline_mode():
    run_migrations_offline()
else:
    run_migrations_online_sync()
