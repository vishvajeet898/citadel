package cache

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/Orange-Health/citadel/conf"
)

type Cache struct {
	Client *redis.Client
}

var cacheObj Cache

var (
	Rdb       *Cache
	keyPrefix = conf.GetConfig().GetString("service_name")
)

func Initialize(ctx context.Context) {
	redisAddr := conf.GetConfig().GetString("redis.address")
	redisPass := conf.GetConfig().GetString("redis.password")
	redisDb := conf.GetConfig().GetInt("redis.db")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       redisDb,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(errors.New(err.Error()))
	}
	cacheObj.Client = client
}

func GetCacheInstance() *Cache {
	return &cacheObj
}

func InitializeCache() CacheLayer {
	return &Cache{
		Client: GetCacheInstance().Client,
	}
}
