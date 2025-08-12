package db

import (
	"app/conf"
	"app/log"
	"errors"
	sqlite "github.com/glebarez/go-sqlite"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"path/filepath"
	"strings"
)

func InitializeSqlite() {
	if !strings.Contains(conf.DB.Type, "sqlite") {
		log.Info("Sqlite dont Enable")
		return
	}
	dsn := conf.DB.DSN
	if dsn == "" {
		dsn = filepath.Join(conf.RootPath, "app.db")
	}

	driverName := "sqlite3WithHooks"
	registerHooks(driverName, &sqlite.Driver{})

	db := getDBConnection(driverName, dsn)
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	DB = db
}

func MigrateSqlite() {
	if !strings.Contains(conf.DB.Type, "sqlite") {
		return
	}

	if !conf.DB.EnableMigrate {
		log.Info("Sqlite dont Enable Migrate")
	}

	sourceDriver, err := iofs.New(MigrationFS, "migrate_sqlite")
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
	} else {
		log.Info("数据库迁移成功！")
	}
}
