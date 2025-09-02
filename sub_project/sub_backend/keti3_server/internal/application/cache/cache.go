package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

var (
	CacheClient *Cache
	configKey   = "cache"
	ConfigError = errors.New("cache config lost")
)

type cacheConfig struct {
	CacheIsOpen  bool          `mapstructure:"cache_is_open"`
	CacheTimeout time.Duration `mapstructure:"cache_timeout"`
	RedisPre     string        `mapstructure:"redis_pre"`
	RedisTimeout time.Duration `mapstructure:"redis_timeout"`
}

type Cache struct {
	*pontos.RedisClient
	conf     map[string]cacheConfig
	cacheMap map[string]*sync.Map
}

type cacheDataStruct struct {
	data      interface{}
	timestamp time.Time
}

func (p *cacheDataStruct) checkTimeout(duration time.Duration) bool {
	return time.Now().Sub(p.timestamp) > duration
}

// NewCache new cache instance
func NewCache() (*Cache, error) {
	if CacheClient != nil {
		return CacheClient, nil
	}
	confMap := map[string]cacheConfig{}
	if err := pontos.Config.UnmarshalKey(configKey, &confMap); err != nil {
		return nil, err
	}

	CacheClient = &Cache{
		RedisClient: pontos.Redis,
		conf:        confMap,
	}
	CacheClient.initCacheMap()
	return CacheClient, nil
}

func (p *Cache) initCacheMap() {
	p.cacheMap = make(map[string]*sync.Map, len(p.conf))
	for key, _ := range p.conf {
		p.cacheMap[key] = new(sync.Map)
	}
}

func (p *Cache) getCache(group, key string) (interface{}, bool) {
	config, ok := p.conf[group]
	if !ok {
		return nil, false
	}

	if data, ok := p.cacheMap[group].Load(key); ok {
		if res, ok := data.(*cacheDataStruct); ok && !res.checkTimeout(config.CacheTimeout) {
			return res.data, true
		}
	}
	return nil, false
}

// HGetCache hash get cache
func (p *Cache) HGetCache(ctx *gin.Context, group, key, field string) (interface{}, bool, error) {
	config, ok := p.conf[group]
	if !ok {
		return nil, false, ConfigError
	}
	key = p.getKey(group, key)
	if config.CacheIsOpen {
		if cacheData, ok := p.getCache(group, field); ok {
			return cacheData, true, nil
		} else {
			result, ok, err := p.hGetRedisCache(ctx, key, field)
			if ok {
				p.setCache(group, key, result)
			}
			return result, ok, err
		}
	} else {
		return p.hGetRedisCache(ctx, key, field)
	}
}

func (p *Cache) hGetRedisCache(ctx *gin.Context, key, field string) (interface{}, bool, error) {
	resultStr, err := p.HGet(ctx, key, field).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, false, err
	} else if err != nil && err.Error() == "redis: nil" {
		return nil, false, nil
	} else {
		var result interface{}
		jsonDecoder := json.NewDecoder(bytes.NewBuffer(resultStr))
		jsonDecoder.UseNumber()
		if err = jsonDecoder.Decode(&result); err != nil {
			return nil, false, err
		} else {
			return result, true, nil
		}
	}
}

// HSetCache hash set cache
func (p *Cache) HSetCache(ctx *gin.Context, group, key, field string, data interface{}) error {
	config, ok := p.conf[group]
	if !ok {
		return ConfigError
	}
	key = p.getKey(group, key)
	if config.CacheIsOpen {
		p.setCache(group, field, data)
	}
	return p.hSetRedisCache(ctx, key, field, data)
}

func (p *Cache) hSetRedisCache(ctx *gin.Context, key, field string, data interface{}) error {
	res, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return p.HSet(ctx, key, res).Err()
}

// GetCache get cache
func (p *Cache) GetCache(ctx *gin.Context, group, key string) (interface{}, bool, error) {
	config, ok := p.conf[group]
	if !ok {
		return nil, false, ConfigError
	}
	key = p.getKey(group, key)
	if config.CacheIsOpen {
		if cacheData, ok := p.getCache(group, key); ok {
			return cacheData, true, nil
		} else {
			data, ok, err := p.getRedisCache(ctx, key)
			if ok {
				p.setCache(group, key, data)
			}
			return data, ok, err
		}
	} else {
		return p.getRedisCache(ctx, key)
	}
}

func (p *Cache) getRedisCache(ctx *gin.Context, key string) (interface{}, bool, error) {
	resultStr, err := p.Get(ctx, key).Bytes()
	if err != nil && err.Error() != "redis: nil" {
		return nil, false, err
	} else if err != nil && err.Error() == "redis: nil" {
		return nil, false, nil
	}

	var result interface{}
	jsonDecoder := json.NewDecoder(bytes.NewBuffer(resultStr))
	jsonDecoder.UseNumber()
	if err = jsonDecoder.Decode(&result); err != nil {
		return nil, false, err
	}
	return result, true, nil
}

// SetCache set cache
func (p *Cache) SetCache(ctx *gin.Context, group, key string, data interface{}) error {
	config, ok := p.conf[group]
	if !ok {
		return ConfigError
	}
	key = p.getKey(group, key)
	if config.CacheIsOpen {
		p.setCache(group, key, data)
	}
	return p.setRedisCache(ctx, key, data, config.RedisTimeout)
}

// DelCache delete cache
func (p *Cache) DelCache(ctx *gin.Context, group, key string) error {
	config, ok := p.conf[group]
	if !ok {
		return ConfigError
	}
	key = p.getKey(group, key)
	if config.CacheIsOpen {
		p.delCache(group, key)
	}
	return p.delRedisCache(ctx, key)
}

func (p *Cache) delCache(group, key string) {
	p.cacheMap[group].Delete(key)
}

func (p *Cache) setCache(group, key string, data interface{}) {
	cData := &cacheDataStruct{
		data:      data,
		timestamp: time.Now(),
	}
	p.cacheMap[group].Store(key, cData)
}

func (p *Cache) setRedisCache(ctx *gin.Context, key string, data interface{}, expiration time.Duration) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.Set(ctx, key, res, expiration).Err()
}

func (p *Cache) delRedisCache(ctx *gin.Context, key string) error {
	return p.Del(ctx, key).Err()
}

func (p *Cache) getKey(group, key string) string {
	return fmt.Sprintf("%s%s", p.conf[group].RedisPre, key)
}
