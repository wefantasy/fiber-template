package db

import (
	"app/conf"
	"app/logger"
	"app/util"
	"context"
	"testing"
	"time"
)

func init() {
	conf.Initialize()
	logger.Initialize()
}

func Test_Redis(t *testing.T) {
	InitializeRedis()
	util.RDB.Set(context.Background(), "test1", "value1", time.Second*30)
	val, err := util.RDB.Get(context.Background(), "test1").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func Test_Sqlite(t *testing.T) {
	InitializeSqlite()
	err := util.DB.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Mysql(t *testing.T) {
	InitializeMysql()
	err := util.DB.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
