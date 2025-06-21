package db

import (
	"app/conf"
	"app/log"
	"database/sql"
	"errors"
	sqlite "github.com/glebarez/go-sqlite"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/qustavo/sqlhooks/v2"
	"path/filepath"
	"strings"
)

func InitializeSqlite() {
	if !strings.Contains(conf.Server.DBType, "sqlite") {
		log.Info("Sqlite dont Enable")
		return
	}
	dbPath := conf.Sqlite.Path
	if dbPath == "" {
		dbPath = filepath.Join(conf.Base.RootPath, "app.db")
	}

	driverName := "sqlite3WithHooks"
	sql.Register(driverName, sqlhooks.Wrap(&sqlite.Driver{}, &Hooks{}))
	db := getDBConnection(driverName, dbPath)
	err := db.Ping()
	if err != nil {
		log.Panic(err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	DB = db
}

func MigrateSqlite() {
	if !strings.Contains(conf.Server.DBType, "sqlite") {
		return
	}

	if !conf.Server.EnableMigrate {
		log.Info("Sqlite dont Enable Migrate")
	}

	sourceDriver, err := iofs.New(MigrationFS, "migrate-sqlite")
	if err != nil {
		log.Fatalf("无法创建迁移源驱动: %v", err)
	}
	dbDriver, err := sqlite3.WithInstance(DB.DB, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("无法创建数据库迁移驱动: %v", err)
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", dbDriver)
	if err != nil {
		log.Fatalf("无法创建 migrate 实例: %v", err)
	}
	log.Info("开始执行数据库迁移...")
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("数据库已经是最新版本，无需迁移。")
		} else {
			log.Fatalf("数据库迁移失败: %v", err)
		}
	}
	log.Info("数据库迁移成功！")
}
