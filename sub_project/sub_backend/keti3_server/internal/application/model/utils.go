package model

import "gorm.io/gorm"

type WithOptions func(*gorm.DB) *gorm.DB

// SetOrder 设置排序
// Usage: model.SetOrder("id desc")
func SetOrder(o interface{}) WithOptions {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Order(o)
		return db
	}
}

// AndWhere 增加查询条件，支持处理<,>,IN 等
// Usage: model.AndWhere("id > ?" ,100)
func AndWhere(query interface{}, args ...interface{}) WithOptions {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(query, args...)
		return db
	}
}

// WithPage 增加分页
func WithPage(page, pageSize int) WithOptions {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize

		if pageSize > 0 {
			db = db.Offset(offset).Limit(pageSize)
		}

		return db
	}
}

// WithPreload 增加预加载
func WithPreload(load string, args ...interface{}) WithOptions {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload(load, args...)
		return db
	}
}
