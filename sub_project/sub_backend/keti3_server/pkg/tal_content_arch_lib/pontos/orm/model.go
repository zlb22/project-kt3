package orm

import (
	"context"

	"gorm.io/gorm"
)

type Model struct {
	TableName string
	DB        *gorm.DB
}

func (m *Model) FindOne(ctx context.Context, fields string, where map[string]interface{}) (map[string]interface{}, error) {
	db := m.DB.Session(&gorm.Session{NewDB: true, Logger: SQLLogger}).Table(m.TableName).WithContext(ctx)
	if fields != "" {
		db = db.Select(fields)
	}
	if len(where) > 0 {
		db = db.Where(where)
	}
	result := map[string]interface{}{}
	db = db.Take(&result).Limit(1)
	return result, db.Error
}

func (m *Model) FindAll(ctx context.Context, fields string, where map[string]interface{}, order string, limit, offset int) ([]map[string]interface{}, error) {
	db := m.DB.Session(&gorm.Session{NewDB: true, Logger: SQLLogger}).Table(m.TableName).WithContext(ctx)
	if fields != "" {
		db = db.Select(fields)
	}
	if len(where) > 0 {
		db = db.Where(where)
	}
	if order != "" {
		db = db.Order(order)
	}
	if limit != 0 {
		db = db.Limit(limit)
	}
	if offset != 0 {
		db = db.Limit(offset)
	}
	result := []map[string]interface{}{}
	db = db.Find(&result)
	return result, db.Error
}

func (m *Model) Count(ctx context.Context, where map[string]interface{}) (int64, error) {
	db := m.DB.Session(&gorm.Session{NewDB: true, Logger: SQLLogger}).Table(m.TableName).WithContext(ctx)
	if len(where) > 0 {
		db = db.Where(where)
	}
	var count int64
	db = db.Count(&count)
	return count, db.Error
}
