package redispool

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

var (
	RedisClient *redis.Pool
	REDIS_HOST string
)


func Init() {
	// 从配置文件获取redis的ip以及db
	REDIS_HOST = "127.0.0.1:6379"
	// 建立连接池
	RedisClient = &redis.Pool{
		MaxIdle:     1,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}
