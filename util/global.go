package util

import (
	"app/code"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap/zapcore"
)

var (
	DB   *sqlx.DB
	RDB  *RedisDB
	Conf *Config
)

type Config struct {
	AppName  string     `toml:"appName"`
	RootPath string     `toml:"rootPath"`
	Server   Server     `toml:"server"`
	Logger   LoggerConf `toml:"logger"`
	Redis    RedisConf  `toml:"redis"`
	Mysql    MysqlConf  `toml:"mysql"`
	Sqlite   SqliteConf `toml:"sqlite"`
}

type Server struct {
	Mode          string `toml:"mode"`
	Address       string `toml:"address"`
	Port          string `toml:"port"`
	EnableMigrate bool   `toml:"enableMigrate"`
	Secret        string `toml:"secret"`
	DBType        string `toml:"dbType"`
}

type LoggerConf struct {
	Level           zapcore.Level `toml:"level"`
	StackTraceLevel zapcore.Level `toml:"stackTraceLevel"`
	Filename        string        `toml:"filename"`
	MaxSize         int           `toml:"maxSize"`
	MaxBackups      int           `toml:"maxBackups"`
	MaxAge          int           `toml:"maxAge"`
	EnableCompress  bool          `toml:"enableCompress"`
}

type MysqlConf struct {
	DSN string `toml:"dsn"`
}

type SqliteConf struct {
	Path string `toml:"path"`
}

type RedisConf struct {
	Enable bool   `toml:"Enable"`
	DSN    string `toml:"dsn"`
	Expire int    `toml:"expire"`
}

type RedisDB struct {
	*redis.Client
}

func (rdb *RedisDB) SetStruct(key string, val interface{}) code.Error {
	if !Conf.Redis.Enable {
		return code.Nil
	}

	jsonData, err := json.Marshal(val)
	if err != nil {
		log.Error(err)
		return code.JsonMarshalFailed
	}
	err = rdb.Client.Set(context.Background(), key, jsonData, time.Duration(Conf.Redis.Expire)*time.Second).Err()
	if err != nil {
		log.Error(err)
		return code.RedisSetDataFailed
	}
	return code.Nil
}

func (rdb *RedisDB) SetStructWithExpire(key string, val interface{}, expire time.Duration) code.Error {
	if !Conf.Redis.Enable {
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
	if !Conf.Redis.Enable {
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
	if !Conf.Redis.Enable {
		return code.Nil
	}
	err := rdb.Client.Del(context.Background(), key).Err()
	if err != nil {
		log.Error(err)
		return code.RedisDeleteDataFailed
	}
	return code.Nil
}
