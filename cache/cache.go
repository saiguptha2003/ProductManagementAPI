package cache

import (
    "context"
    "github.com/go-redis/redis/v8"
    "log"
    "time"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() {
    rdb = redis.NewClient(&redis.Options{
        Addr: "localhost:6379", 
    })
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
}
func Set(key string, value string, expiration time.Duration) error {
    return rdb.Set(ctx, key, value, expiration).Err()
}
func Get(key string) (string, error) {
    return rdb.Get(ctx, key).Result()
}

func Delete(key string) error {
    return rdb.Del(ctx, key).Err()
}
