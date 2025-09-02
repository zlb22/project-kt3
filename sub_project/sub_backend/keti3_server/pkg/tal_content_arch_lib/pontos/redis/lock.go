package redis

import (
	"context"
	"errors"
	"time"
)

func (c *Client) GetLock(ctx context.Context, key string, ttl, timeout time.Duration) (bool, error) {
	now := time.Now()
	for {
		ok := c.SetNX(ctx, key, "1", ttl).Val()
		if ok {
			return true, nil
		}

		time.Sleep(time.Millisecond * 100)
		if time.Now().Sub(now) > timeout {
			return false, errors.New("get lock timeout")
		}
	}
}

func (c *Client) Unlock(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}
