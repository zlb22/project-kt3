from sqlalchemy import BigInteger, Column, String, TIMESTAMP, DATE, func
from ..db.base import Base


class StudentSubmit(Base):
    __tablename__ = "keti3_student_submit"

    id = Column(BigInteger().with_variant(BigInteger, "mysql"), primary_key=True, autoincrement=True)
    submit_id = Column(BigInteger, nullable=True)
    uid = Column(BigInteger, nullable=False, index=True)
    school_id = Column(BigInteger, nullable=False)
    student_name = Column(String(64), nullable=False)
    student_num = Column(String(64), nullable=False)
    grade_id = Column(BigInteger, nullable=False)
    class_name = Column(String(64), nullable=False)
    voice_url = Column(String(255), nullable=True)
    oplog_url = Column(String(255), nullable=True)
    screenshot_url = Column(String(255), nullable=True)
    date = Column(DATE, nullable=True)
    created_at = Column(TIMESTAMP, nullable=False, server_default=func.current_timestamp())
    updated_at = Column(TIMESTAMP, nullable=False, server_default=func.current_timestamp(), onupdate=func.current_timestamp())
