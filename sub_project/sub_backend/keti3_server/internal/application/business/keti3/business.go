package keti3

import (
	"context"

	"keti3/internal/application/model"
	"keti3/internal/application/service/base_video"
)

type Business interface {
	// 获取上传 oss 配置
	GetSTSAuth(ctx context.Context) (ret GetSTSAuthRet, err error)
	// 获取自增 ID
	GetUIDByIncrease(ctx context.Context) (ret GetUIDByIncreaseRet, err error)
	// 保存日志
	SaveLog(ctx context.Context, logParam SaveLogParam) error
	// 导出日志
	ExportLog(ctx context.Context) (string, error)

	// 获取配置列表
	ConfigList(ctx context.Context) (ret ConfigListRet, err error)
	// 学生登陆
	StudentLogin(ctx context.Context, param StudentLoginParam) (ret StudentLoginRet, err error)
	// 老师登陆
	TeacherLogin(ctx context.Context, param TeacherLoginParam) (ret TeacherLoginRet, err error)
	// 检查老师登陆
	CheckTeacherLogin(ctx context.Context, token string) error
	// 日志列表
	LogList(ctx context.Context, param LogListParam) (ret LogListRet, err error)
}

type business struct {
	logModel      *model.Keti3OpLogModel
	baseVideoSrv  *base_video.Client
	stuModel      *model.Keti3StudentInfoModel
	stuSumitModel *model.Keti3StudentSubmitModel
}

// NewBusiness return the instance of interface
func NewBusiness() Business {
	return &business{
		logModel:      model.NewKeti3OpLogModel(),
		baseVideoSrv:  base_video.NewClient(),
		stuModel:      model.NewKeti3StudentInfoModel(),
		stuSumitModel: model.NewKeti3StudentSubmitModel(),
	}
}

var tokenSalt = "keti3"
