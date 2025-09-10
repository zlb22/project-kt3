from sqlalchemy import create_engine, Column, Integer, String, DateTime, Text, BigInteger
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from datetime import datetime
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# MySQL Configuration
MYSQL_DSN = os.getenv("MYSQL_DSN", "mysql+pymysql://root:123456@localhost:3306/keti3")

# SQLAlchemy setup
engine = create_engine(MYSQL_DSN, echo=False)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

class Student(Base):
    __tablename__ = "students"
    
    id = Column(Integer, primary_key=True, index=True)
    username = Column(String(50), unique=True, index=True, nullable=False)
    school = Column(String(100), nullable=False)
    grade = Column(String(20), nullable=False)
    password = Column(String(128), nullable=True)
    is_active = Column(Integer, nullable=False, default=1)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=True)

class OperationLog(Base):
    __tablename__ = "operation_logs"
    
    id = Column(Integer, primary_key=True, index=True)
    uid = Column(Integer, nullable=False)
    action = Column(String(100), nullable=False)
    details = Column(Text)
    timestamp = Column(DateTime, default=datetime.utcnow)

class TwentyFourRecord(Base):
    __tablename__ = "twentyfour_records"

    id = Column(BigInteger, primary_key=True, autoincrement=True)
    uid = Column(Integer, nullable=False, index=True)
    payload = Column(Text, nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)

# Create tables
def create_tables():
    Base.metadata.create_all(bind=engine)

# Database dependency
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
