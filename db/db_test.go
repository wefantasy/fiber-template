package db

import (
	"app/conf"
	"app/log"
	"context"
	"testing"
	"time"
)

func TestMigrateSqlite(t *testing.T) {
	conf.Initialize()
	log.Initialize()
	InitializeSqlite()
	MigrateSqlite()
}

func TestMigrateMysql(t *testing.T) {
	InitializeMysql()
	MigrateMysql()
}

func Test_Redis(t *testing.T) {
	InitializeRedis()
	RDB.Set(context.Background(), "test1", "value1", time.Second*30)
	val, err := RDB.Get(context.Background(), "test1").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func Test_Sqlite(t *testing.T) {
	InitializeSqlite()
	err := DB.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Mysql(t *testing.T) {
	InitializeMysql()
	err := DB.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
