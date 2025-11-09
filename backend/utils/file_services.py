import hashlib
from sqlmodel import select

def generate_checksum(file_name: str) -> str:
    with open(file_name, 'rb') as file:
        file_hash = hashlib.file_digest(file, 'sha256').hexdigest()
    return file_hash

def verify_checksum(file_name: str) -> bool:
    file_hash = generate_checksum(file_name)
    statement = select(Blog.file_hash)where(Blog.file_name == file_name)

    db_file_hash = session.exec(statement).first()

    if db_file_hash == file_hash:
        return True
    
    return False