package model

import (
	"context"

	"keti3/internal/application/schema"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"gorm.io/gorm"
)

type Keti3StudentInfoModel struct {
	Orm *gorm.DB
}

func NewKeti3StudentInfoModel() *Keti3StudentInfoModel {
	return &Keti3StudentInfoModel{
		Orm: pontos.MysqlMap["base_keti"].Orm,
	}
}

func (impl *Keti3StudentInfoModel) db(ctx context.Context) *gorm.DB {
	return impl.Orm.WithContext(ctx)
}

func (impl *Keti3StudentInfoModel) Create(ctx context.Context, data *schema.Keti3StudentInfo) (int64, error) {
	err := impl.db(ctx).Create(data).Error
	return data.ID, err
}

func (impl *Keti3StudentInfoModel) GetStudent(ctx context.Context, param *schema.Keti3StudentInfo) (*schema.Keti3StudentInfo, error) {
	var data *schema.Keti3StudentInfo
	err := impl.db(ctx).Where(param).First(&data).Error
	return data, err
}
