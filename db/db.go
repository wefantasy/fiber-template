package db

import (
	"app/log"
	"context"
	"database/sql"
	"embed"
	"github.com/jmoiron/sqlx"
	"time"
)

func Initialize() {
	InitializeRedis()
	InitializeSqlite()
	InitializeMysql()
}

func Migrate() {
	MigrateSqlite()
	MigrateMysql()
}

//go:embed migrate-sqlite/*.sql
//go:embed migrate-mysql/*.sql
var MigrationFS embed.FS

var DB *sqlx.DB
var RDB *RedisDB

type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	log.Infof("Exec SQL: %s %v", query, args)
	return context.WithValue(ctx, "beginTime", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("beginTime").(time.Time)
	log.Infof("Above SQL Used Time: %s", time.Since(begin))
	return ctx, nil
}

func getDBConnection(driverName, dataSourceName string) *sqlx.DB {
	sqlDB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	db := sqlx.NewDb(sqlDB, driverName)
	return db
}
