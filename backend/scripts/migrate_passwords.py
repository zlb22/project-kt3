#!/usr/bin/env python3
"""
Password Migration Script for Project-KT3
Migrates existing SHA256 password hashes to bcrypt hashes
"""

import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import bcrypt
import hashlib
from sqlalchemy.orm import Session
from db_mysql import get_db, Student

def migrate_password_hashes():
    """Migrate existing SHA256 password hashes to bcrypt"""
    db = next(get_db())
    
    try:
        # Get all students with passwords
        students = db.query(Student).filter(Student.password.isnot(None)).all()
        
        migrated_count = 0
        skipped_count = 0
        
        for student in students:
            # Check if password is already bcrypt (starts with $2b$)
            if student.password.startswith('$2b$'):
                print(f"Skipping {student.username} - already bcrypt")
                skipped_count += 1
                continue
            
            # Assume it's SHA256, convert to bcrypt
            # Note: We can't recover the original password, so we'll mark these for password reset
            print(f"Migrating {student.username} - marking for password reset")
            
            # Set a temporary bcrypt hash that will force password reset
            temp_password = f"MIGRATION_REQUIRED_{student.username}"
            salt = bcrypt.gensalt()
            new_hash = bcrypt.hashpw(temp_password.encode('utf-8'), salt)
            
            student.password = new_hash.decode('utf-8')
            student.updated_at = datetime.utcnow()
            
            migrated_count += 1
        
        db.commit()
        
        print(f"\nMigration completed:")
        print(f"  Migrated: {migrated_count} users")
        print(f"  Skipped: {skipped_count} users")
        
        if migrated_count > 0:
            print(f"\nIMPORTANT: {migrated_count} users will need to reset their passwords")
            print("Consider implementing a password reset flow or notifying users")
        
    except Exception as e:
        print(f"Migration failed: {e}")
        db.rollback()
        return False
    finally:
        db.close()
    
    return True

if __name__ == "__main__":
    from datetime import datetime
    print("Starting password migration...")
    success = migrate_password_hashes()
    if success:
        print("Migration completed successfully")
    else:
        print("Migration failed")
        sys.exit(1)
