# ruff: noqa: F401
from sqlmodel import create_engine, Session, SQLModel
from models.schema import BlogTagLink, Blog, BlogContent, Tag, SubscriberEmail
import os

# Database URL from environment variables
DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://postgres@localhost:5450/blog")

# Create engine with connection pooling configuration
engine = create_engine(
    DATABASE_URL, 
    echo=True,
    pool_size=10,
    max_overflow=20,
    pool_pre_ping=True,  # Validates connections before use
    pool_recycle=3600    # Recycle connections every hour
)

def create_db_and_tables():
    """Create database tables"""
    SQLModel.metadata.create_all(engine)

def get_session():
    """Get database session with error handling"""
    try:
        with Session(engine) as session:
            yield session
    except Exception as e:
        print(f"Database session error: {e}")
        raise