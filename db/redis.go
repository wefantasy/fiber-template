package db

import (
	"app/code"
	"app/conf"
	"app/log"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var RDB *RedisDB

func InitializeRedis() {
	if !conf.Redis.Enable {
		log.Info("Redis dont Enable")
		return
	}

	opt, _ := redis.ParseURL(conf.Redis.DSN)
	rdb := redis.NewClient(opt)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Panic(err)
	}
	RDB = &RedisDB{
		Client: rdb,
	}
}

type RedisDB struct {
	*redis.Client
}

func (rdb *RedisDB) SetStruct(key string, val interface{}) error {
	if !conf.Redis.Enable {
		return nil
	}

	jsonData, err := json.Marshal(val)
	if err != nil {
		log.Error(err)
		return err
	}
	err = rdb.Client.Set(context.Background(), key, jsonData, time.Duration(conf.Redis.Expire)*time.Second).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (rdb *RedisDB) SetStructWithExpire(key string, val interface{}, expire time.Duration) error {
	if !conf.Redis.Enable {
		return nil
	}

	jsonData, err := json.Marshal(val)
	if err != nil {
		log.Error(err)
		return err
	}

	err = rdb.Client.Set(context.Background(), key, jsonData, expire).Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (rdb *RedisDB) GetStruct(key string, obj interface{}) error {
	if !conf.Redis.Enable {
		return nil
	}

	jsonData, err := rdb.Client.Get(context.Background(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return code.RedisKeyNotExist
		}
		log.Error(err)
		return nil
	}

	err = json.Unmarshal(jsonData, obj)
	if err != nil {
		log.Error(err)
		return nil
	}
	return nil
}

func (rdb *RedisDB) Delete(key string) error {
	if !conf.Redis.Enable {
		return nil
	}
	err := rdb.Client.Del(context.Background(), key).Err()
	if err != nil {
		log.Error(err)
		return nil
	}
	return nil
}
