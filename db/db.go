package db

import (
	"app/conf"
	"app/log"
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"github.com/jmoiron/sqlx"
	"github.com/qustavo/sqlhooks/v2"
	"time"
)

func Initialize() {
	InitializeRedis()
	InitializeSqlite()
	InitializeMysql()
	if conf.DB.EnableMigrate {
		Migrate()
	}
}

func Migrate() {
	MigrateSqlite()
	MigrateMysql()
}

//go:embed migrate_sqlite/*.sql
//go:embed migrate_mysql/*.sql
var MigrationFS embed.FS

var DB *sqlx.DB
var RDB *RedisDB

type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	log.Infof("Exec SQL: \n%s %v", query, args)
	return context.WithValue(ctx, "beginTime", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {
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
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	return db
}

func registerHooks(driverName string, driver driver.Driver) {
	var driverIsRegistered bool
	for _, d := range sql.Drivers() {
		if d == driverName {
			driverIsRegistered = true
			break
		}
	}
	if !driverIsRegistered {
		sql.Register(driverName, sqlhooks.Wrap(driver, &Hooks{}))
	}
}
