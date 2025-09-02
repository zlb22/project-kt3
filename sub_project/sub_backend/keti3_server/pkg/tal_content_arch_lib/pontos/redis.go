package pontos

import (
	"keti3/pkg/tal_content_arch_lib/pontos/bootstrap"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"
	"keti3/pkg/tal_content_arch_lib/pontos/redis"
)

type RedisClient struct {
	*redis.Client
}

func InitRedis(logger *logger.ZLogger) error {
	if Config.IsSet(redisSection) {
		conf := &redis.RedisConfig{}
		if err := Config.UnmarshalKey(redisSection, conf); err != nil {
			return err
		}

		client, err := redis.NewClient(conf)
		if err != nil {
			return err
		}
		Redis = &RedisClient{
			client,
		}
		Redis.AddHook(redis.NewLoggerHook(logger))
		Redis.AddHook(redis.NewOTELHook(conf.DB))
		return nil
	}
	return nil
}

func CloseRedis() bootstrap.AfterServerFunc {
	return func() {
		if Config.IsSet(redisSection) {
			_ = Redis.Close()
		}
	}
}

func InitRedisMap(logger *logger.ZLogger) error {
	if Config.IsSet(redisMapSection) {
		confMap := make(map[string]*redis.RedisConfig)
		if err := Config.UnmarshalKey(redisMapSection, &confMap); err != nil {
			return err
		}

		for name, config := range confMap {
			if client, err := redis.NewClient(config); err != nil {
				return err
			} else {
				RedisMap[name] = &RedisClient{client}
				RedisMap[name].AddHook(redis.NewLoggerHook(logger))
				RedisMap[name].AddHook(redis.NewOTELHook(config.DB))
			}
		}
	}

	return nil
}

func CloseRedisMap() bootstrap.AfterServerFunc {
	return func() {
		if Config.IsSet(redisMapSection) {
			for _, client := range RedisMap {
				_ = client.Close()
			}
		}
	}
}
