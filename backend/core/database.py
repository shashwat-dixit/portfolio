from sqlmodel import create_engine, Session
import os

# Database URL from environment variables
DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://user:password@localhost:5432/portfolio")

# Create engine
engine = create_engine(DATABASE_URL, echo=True)

def get_session():
    """Get database session"""
    with Session(engine) as session:
        yield session