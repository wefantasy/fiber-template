package db

import (
	"app/conf"
	"app/log"
	"app/model"
	"database/sql"
	sqlite "github.com/glebarez/go-sqlite"
	"github.com/qustavo/sqlhooks/v2"
	"os"
	"path/filepath"
	"strings"
)

func Initialize() {
	InitializeRedis()
	InitializeSqlite()
	InitializeMysql()
}

func InitializeSqlite() {
	if !strings.Contains(conf.Server.DBType, "sqlite") {
		log.Info("Sqlite dont Enable")
		return
	}
	dbPath := conf.Sqlite.Path
	if dbPath == "" {
		dbPath = filepath.Join(conf.Base.RootPath, "app.db")
	}

	initSqlStr := ""
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Info("Sqlite数据库不存在，需要初始化...")
		sqlBytes, err := model.SqlFS.ReadFile("sqlite.sql")
		if err != nil {
			log.Panic(err)
		}
		initSqlStr = string(sqlBytes)
	}

	driverName := "sqlite3WithHooks"
	sql.Register(driverName, sqlhooks.Wrap(&sqlite.Driver{}, &Hooks{}))
	db := getDBConnection(driverName, dbPath)
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
