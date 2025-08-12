from pydantic import BaseModel
from datetime import datetime
from uuid import UUID

class Blog(BaseModel):
    id: UUID
    title: str
    modified_at: datetime
    slug: str