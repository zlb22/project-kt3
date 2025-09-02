package keti3

import (
	"context"

	"keti3/internal/application/constant"

	"keti3/pkg/tal_content_arch_lib/pontos"
)

type GetUIDByIncreaseRet struct {
	UID int64 `json:"uid"`
}

func (b business) GetUIDByIncrease(ctx context.Context) (ret GetUIDByIncreaseRet, err error) {
	id, err := getIncreaseIDByRedis(ctx)
	if err != nil {
		return
	}
	ret.UID = id
	return
}

// redis 自增id
func getIncreaseIDByRedis(ctx context.Context) (int64, error) {
	id, err := pontos.Redis.Incr(ctx, constant.IncreaseUIDKey).Result()
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("getIncreaseIDByRedis err:%v", err)
	}
	pontos.Logger.WithContext(ctx).Infof("getIncreaseIDByRedis id:%d", id)

	return id, err
}
