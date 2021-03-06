package redisfx

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"

	"github.com/nkhang/pluto/pkg/cache"
)

func provideRedisClient() (redis.UniversalClient, error) {
	addr := viper.GetString("redis.url")
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	err := client.Ping().Err()
	return client, err
}

func provideCacheRepository(client redis.UniversalClient) cache.Cache {
	client.FlushAll()
	return cache.New(client, cache.WithExpireTime(time.Hour))
}
