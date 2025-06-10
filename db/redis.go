package db

import (
	"app/code"
	"app/conf"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
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

func (rdb *RedisDB) SetStruct(key string, val interface{}) code.Error {
	if !conf.Redis.Enable {
		return code.Nil
	}

	jsonData, err := json.Marshal(val)
	if err != nil {
		log.Error(err)
		return code.JsonMarshalFailed
	}
	err = rdb.Client.Set(context.Background(), key, jsonData, time.Duration(conf.Redis.Expire)*time.Second).Err()
	if err != nil {
		log.Error(err)
		return code.RedisSetDataFailed
	}
	return code.Nil
}

func (rdb *RedisDB) SetStructWithExpire(key string, val interface{}, expire time.Duration) code.Error {
	if !conf.Redis.Enable {
		return code.Nil
	}

	jsonData, err := json.Marshal(val)
	if err != nil {
		log.Error(err)
		return code.JsonMarshalFailed
	}

	err = rdb.Client.Set(context.Background(), key, jsonData, expire).Err()
	if err != nil {
		log.Error(err)
		return code.RedisSetDataFailed
	}
	return code.Nil
}

func (rdb *RedisDB) GetStruct(key string, obj interface{}) code.Error {
	if !conf.Redis.Enable {
		return code.Nil
	}

	jsonData, err := rdb.Client.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return code.RedisKeyNotExist
		}
		log.Error(err)
		return code.RedisGetDataFailed
	}

	err = json.Unmarshal(jsonData, obj)
	if err != nil {
		log.Error(err)
		return code.JsonUnmarshalFailed
	}
	return code.Nil
}

func (rdb *RedisDB) Delete(key string) code.Error {
	if !conf.Redis.Enable {
		return code.Nil
	}
	err := rdb.Client.Del(context.Background(), key).Err()
	if err != nil {
		log.Error(err)
		return code.RedisDeleteDataFailed
	}
	return code.Nil
}
