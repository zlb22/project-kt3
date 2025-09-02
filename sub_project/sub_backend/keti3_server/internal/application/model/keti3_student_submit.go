package model

import (
	"context"

	"keti3/internal/application/schema"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"gorm.io/gorm"
)

type Keti3StudentSubmitModel struct {
	Orm *gorm.DB
	tb  schema.Keti3StudentSubmit
}

func NewKeti3StudentSubmitModel() *Keti3StudentSubmitModel {
	return &Keti3StudentSubmitModel{
		Orm: pontos.MysqlMap["base_keti"].Orm,
		tb:  schema.Keti3StudentSubmit{},
	}
}

func (impl *Keti3StudentSubmitModel) db(ctx context.Context) *gorm.DB {
	return impl.Orm.WithContext(ctx).Table(impl.tb.TableName())
}

func (impl *Keti3StudentSubmitModel) Create(ctx context.Context, data *schema.Keti3StudentSubmit) (int64, error) {
	err := impl.db(ctx).Create(data).Error
	return data.ID, err
}

func (impl *Keti3StudentSubmitModel) GetList(ctx context.Context, options []WithOptions) ([]*schema.Keti3StudentSubmit, error) {
	var data []*schema.Keti3StudentSubmit

	db := impl.db(ctx)
	for _, option := range options {
		db = option(db)
	}
	err := db.Find(&data).Error
	return data, err
}

func (impl *Keti3StudentSubmitModel) Count(ctx context.Context, options []WithOptions) (int64, error) {
	var count int64

	db := impl.db(ctx)
	for _, option := range options {
		db = option(db)
	}
	err := db.Count(&count).Error
	return count, err
}
