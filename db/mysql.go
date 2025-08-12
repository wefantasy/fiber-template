package db

import (
	"app/conf"
	"app/log"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysql2 "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"strings"
)

func InitializeMysql() {
	if !strings.Contains(conf.DB.Type, "mysql") {
		log.Info("Mysql dont Enable")
		return
	}
	driverName := "mysqlWithHooks"
	registerHooks(driverName, &mysql.MySQLDriver{})

	// 不存在则创建
	parsedDSN, err := mysql.ParseDSN(conf.DB.DSN)
	if err != nil {
		log.Panic(err)
	}
	dbName := parsedDSN.DBName
	parsedDSN.DBName = ""
	initDb := getDBConnection(driverName, parsedDSN.FormatDSN())
	defer initDb.Close()
	if err = initDb.Ping(); err != nil {
		log.Panic(err)
	}
	createDbSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName)
	_, err = initDb.Exec(createDbSql)
	if err != nil {
		log.Fatalf("无法连接到 MySQL 服务器: %v", err)
	}

	db := getDBConnection(driverName, conf.DB.DSN)
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	DB = db
}

func MigrateMysql() {
	if !strings.Contains(conf.DB.Type, "mysql") {
		return
	}

	if !conf.DB.EnableMigrate {
		log.Info("Mysql dont Enable Migrate")
	}

	sourceDriver, err := iofs.New(MigrationFS, "migrate_mysql")
	if err != nil {
		log.Fatalf("无法创建迁移源驱动: %v", err)
	}
	dbDriver, err := mysql2.WithInstance(DB.DB, &mysql2.Config{})
	if err != nil {
		log.Fatalf("无法创建数据库迁移驱动: %v", err)
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", dbDriver)
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
