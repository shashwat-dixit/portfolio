import os

class Settings:
    PROJECT_NAME: str = "Portfolio Blog API"
    VERSION: str = "1.0.0"
    DATABASE_URL: str = os.getenv("DATABASE_URL", "postgresql://user:password@localhost:5432/portfolio")
    GITHUB_TOKEN: str = os.getenv("GITHUB_TOKEN", "")
    GITHUB_WEBHOOK_SECRET: str = os.getenv("GITHUB_WEBHOOK_SECRET", "")
    REPO_NAME: str = os.getenv("REPO_NAME", "")
    GITHUB_BRANCH: str = os.getenv("GITHUB_BRANCH", "main")

settings = Settings()