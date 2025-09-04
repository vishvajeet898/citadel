package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
)

type CacheLayer interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, data interface{}, expiryTime int) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	HSet(ctx context.Context, key, field string, value interface{}, expiry time.Duration) error
	HSetAll(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]interface{}, error)
}

// Set: set a key/value
func (c *Cache) Set(ctx context.Context, key string, data interface{}, expiryTime int) error {
	key = keyPrefix + ":" + key
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = c.Client.Set(ctx, key, value, time.Duration(expiryTime)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get: get a key/value
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	key = keyPrefix + ":" + key
	val, err := c.Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		log.Println(commonConstants.ERROR_KEY_DOES_NOT_EXIST, ":", key)
		return err
	} else if err != nil {
		log.Println(commonConstants.ERROR_WHILE_FETCHING_DATA_FROM_REDIS, ", Key: ", key)
		return err
	}
	err = json.Unmarshal(val, dest)
	return err
}

// Exists: check if a key exists
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	key = keyPrefix + ":" + key
	cmd := c.Client.Exists(ctx, key)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return cmd.Val() == 1, nil
}

// Delete: delete a key
func (c *Cache) Delete(ctx context.Context, key string) error {
	key = keyPrefix + ":" + key
	return c.Client.Del(ctx, key).Err()
}

// HSet: set a key/value in hash
func (c *Cache) HSet(ctx context.Context, key, field string, value interface{}, expiry time.Duration) error {
	key = keyPrefix + ":" + key
	byteStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := c.Client.HSet(ctx, key, field, string(byteStr))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return c.ExpireAt(ctx, key, expiry)
}

// HSetAll: set all key/value in hash
func (c *Cache) HSetAll(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	key = keyPrefix + ":" + key
	var valueMap map[string]interface{}
	byteStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err1 := json.Unmarshal(byteStr, &valueMap)
	if err1 != nil {
		return err1
	}
	for k, v := range valueMap {
		bs, err2 := json.Marshal(v)
		if err2 != nil {
			return err2
		}
		cmd := c.Client.HSet(ctx, key, k, string(bs))
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return c.ExpireAt(ctx, key, expiry)
}

// HGet: get a key/value from hash
func (c *Cache) HGet(ctx context.Context, key, field string) (string, error) {
	key = keyPrefix + ":" + key
	value, err := c.Client.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// HGetAll: get all key/value from hash
func (c *Cache) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	key = keyPrefix + ":" + key
	response := make(map[string]interface{})
	res, err := c.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return response, err
	}
	for k, value := range res {
		var subValue interface{}
		e := json.Unmarshal([]byte(value), &subValue)
		if e != nil {
			response[k] = value
		} else {
			response[k] = subValue
		}
	}
	return response, err
}

// ExpireAt: set expiry time for a key
func (c *Cache) ExpireAt(ctx context.Context, key string, expire time.Duration) error {
	cmd := c.Client.Expire(ctx, key, expire)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
