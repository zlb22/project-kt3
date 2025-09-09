from sqlalchemy import BigInteger, Column, String, TIMESTAMP, JSON, func
from ..db.base import Base


class OpLog(Base):
    __tablename__ = "keti3_op_log"

    id = Column(BigInteger().with_variant(BigInteger, "mysql"), primary_key=True, autoincrement=True)
    submit_id = Column(BigInteger, nullable=True)
    uid = Column(BigInteger, nullable=False, index=True)
    op_time = Column(TIMESTAMP, nullable=False, server_default=func.current_timestamp())
    op_type = Column(String(32), nullable=False)
    op_object = Column(String(64), nullable=True)
    object_no = Column(String(64), nullable=True)
    object_name = Column(String(128), nullable=True)
    data_before = Column(JSON, nullable=True)
    data_after = Column(JSON, nullable=True)
    voice_url = Column(String(255), nullable=True)
    screenshot_url = Column(String(255), nullable=True)
    deleted_at = Column(TIMESTAMP, nullable=True)
    created_at = Column(TIMESTAMP, nullable=False, server_default=func.current_timestamp())
    updated_at = Column(TIMESTAMP, nullable=False, server_default=func.current_timestamp(), onupdate=func.current_timestamp())
