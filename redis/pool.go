package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	// Pool contains the Redis thread pool
	Pool *redis.Pool
)

// InitPool initializes the Redis thread pool.
func InitPool() error {
	redisHost := fmt.Sprintf("%s:%s",
		os.Getenv("B_REDIS_HOST"),
		os.Getenv("B_REDIS_PORT"),
	)

	Pool = newPool(redisHost)
	return Ping()
}

// ClosePool closes the Redis thread pool.
func ClosePool() {
	Pool.Close()
}

func newPool(redisHost string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisHost)
		},
	}
}

// Ping pings the connection to Redis and returns an error
// if the connection cannot be established.
func Ping() error {
	conn := Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("PING redis failed: %v", err)
	}
	return nil
}
