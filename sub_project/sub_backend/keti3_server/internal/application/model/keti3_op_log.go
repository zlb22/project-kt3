package model

import (
	"context"

	"keti3/internal/application/schema"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"gorm.io/gorm"
)

type Keti3OpLogModel struct {
	Orm *gorm.DB
}

func NewKeti3OpLogModel() *Keti3OpLogModel {
	return &Keti3OpLogModel{
		Orm: pontos.MysqlMap["base_keti"].Orm,
	}
}

func (impl *Keti3OpLogModel) db(ctx context.Context) *gorm.DB {
	return impl.Orm.WithContext(ctx)
}

func (impl *Keti3OpLogModel) BatchCreate(ctx context.Context, data []schema.Keti3OpLog) error {
	return impl.db(ctx).CreateInBatches(data, len(data)).Error
}

func (impl *Keti3OpLogModel) GetList(ctx context.Context, limit, offset int) ([]schema.Keti3OpLog, error) {
	data := make([]schema.Keti3OpLog, 0)
	err := impl.db(ctx).
		Table(schema.Keti3OpLog{}.TableName()).
		Limit(limit).Offset(offset).
		Find(&data).Error

	return data, err
}

func (impl *Keti3OpLogModel) GetCount(ctx context.Context) (int64, error) {
	var (
		count int64
	)
	err := impl.db(ctx).
		Table(schema.Keti3OpLog{}.TableName()).
		Count(&count).Error

	return count, err
}
