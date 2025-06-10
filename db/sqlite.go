package db

import (
	"app/conf"
	"database/sql"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
	"os"
	"path"
)

func InitializeSqlite() {
	if conf.Server.DBType != "sqlite" {
		log.Info("Sqlite dont Enable")
		return
	}

	initSqlStr := ""
	_, err := os.Stat(conf.Sqlite.Path)
	if os.IsNotExist(err) {
		log.Info("Sqlite数据库不存在，需要初始化...")
		sqlFilePath := path.Join(conf.RootPath, "model", "sqlite.sql")
		sqlBytes, err := os.ReadFile(sqlFilePath)
		if err != nil {
			log.Panic(err)
		}
		initSqlStr = string(sqlBytes)
	}

	driverName := "sqlite3WithHooks"
	sql.Register(driverName, sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, &Hooks{}))
	db := getDBConnection(driverName, conf.Sqlite.Path)
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	if initSqlStr != "" {
		_, err = db.Exec(initSqlStr)
		if err != nil {
			log.Panic(err)
		}
		log.Info("数据库初始化成功")
	}
	DB = db
}
