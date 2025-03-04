package util

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type RedisCache struct {
	pool *redis.Pool
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", ":6379")
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func (rc *RedisCache) Set(key string, value interface{}) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		log.Println("Error setting key in Redis:", err)
		return err
	}

	return nil
}

func (rc *RedisCache) Get(key string) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Println("Error getting key from Redis:", err)
		return nil, err
	}

	return value, nil
}
