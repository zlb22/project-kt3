package keti3

import (
	"context"
	"errors"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/go-redis/redis/v8"
)

func (b business) CheckTeacherLogin(ctx context.Context, token string) error {
	hashCode := b.formatRdsKey(token)
	val, err := pontos.Redis.Get(ctx, hashCode).Result()
	if err != nil && err != redis.Nil {
		pontos.Logger.Error("CheckTeacherLogin", "err", err)
		return err
	}
	if val == "" {
		return errors.New("请登陆")
	}
	return nil
}
