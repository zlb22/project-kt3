from sqlalchemy import BigInteger, Column, String, TIMESTAMP, func
from ..db.base import Base


class StudentInfo(Base):
    __tablename__ = "keti3_student_info"

    id = Column(BigInteger().with_variant(BigInteger, "mysql"), primary_key=True, autoincrement=True)
    uid = Column(BigInteger, nullable=False, index=True)
    school_id = Column(BigInteger, nullable=False, index=True)
    student_name = Column(String(64), nullable=False)
    student_num = Column(String(64), nullable=False)
    grade_id = Column(BigInteger, nullable=False)
    class_name = Column(String(64), nullable=False)
    created_at = Column(TIMESTAMP, nullable=False, server_default=func.current_timestamp())
    updated_at = Column(TIMESTAMP, nullable=False, server_default=func.current_timestamp(), onupdate=func.current_timestamp())
