package redis

import (
	"context"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"

	"github.com/go-redis/redis/v8"
)

type pLogHook struct {
	logger *logger.ZLogger
}

var _ redis.Hook = (*pLogHook)(nil)

func (ph *pLogHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (ph *pLogHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	log := ph.logger.WithContext(ctx)

	if err := cmd.Err(); err != nil {
		if err != redis.Nil {
			log = ph.logger.WithError(err)
		}
	}

	log.WithContext(ctx).WithFields(map[string]interface{}{
		"x_cmd": CmdString(cmd),
		"x_tag": constant.LogRedisTag,
	}).Info(cmd.FullName())

	return nil
}

func (ph *pLogHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (ph *pLogHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	summary, cmdsString := CmdsString(cmds)
	log := ph.logger.WithContext(ctx)

	if err := cmds[0].Err(); err != nil {
		if err != redis.Nil {
			log = ph.logger.WithError(err)
		}
	}

	log.WithContext(ctx).WithFields(map[string]interface{}{
		"x_cmd":     cmdsString,
		"x_num_cmd": summary,
		"x_tag":     constant.LogRedisTag,
	}).Info("pipeline " + summary)
	return nil
}

// NewLoggerHook return a redis Hook for logger
func NewLoggerHook(pLogger *logger.ZLogger) redis.Hook {
	return &pLogHook{
		pLogger,
	}
}
