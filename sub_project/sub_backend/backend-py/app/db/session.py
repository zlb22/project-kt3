import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from dotenv import load_dotenv

# Load .env file
load_dotenv()

MYSQL_DSN = os.getenv("MYSQL_DSN", "")
if not MYSQL_DSN:
    # Example: mysql+pymysql://kt3_user:123456@localhost:3306/kt3?charset=utf8mb4
    raise RuntimeError("MYSQL_DSN not configured in environment")

engine = create_engine(MYSQL_DSN, pool_pre_ping=True, pool_recycle=3600)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
