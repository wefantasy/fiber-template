package db

import (
	"app/conf"
	"app/log"
	"app/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/qustavo/sqlhooks/v2"
	"strings"
	"time"
)

var DB *sqlx.DB

type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	log.Infof("Exec SQL\n\t%s %v", query, args)
	return context.WithValue(ctx, "beginTime", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("beginTime").(time.Time)
	log.Infof("Above SQL Used Time: %s", time.Since(begin))
	return ctx, nil
}

func InitializeMysql() {
	if !strings.Contains(conf.Server.DBType, "mysql") {
		log.Info("Mysql dont Enable")
		return
	}
	driverName := "mysqlWithHooks"
	sql.Register(driverName, sqlhooks.Wrap(&mysql.MySQLDriver{}, &Hooks{}))
	db := getDBConnection(driverName, conf.Mysql.DSN)
	err := db.Ping()
	if err != nil {
		if strings.Contains(err.Error(), "Unknown database") {
			log.Info("Mysql数据库不存在，需要初始化...")
			dbCgf, err := mysql.ParseDSN(conf.Mysql.DSN)
			if err != nil {
				log.Panic(err)
			}
			dbName := dbCgf.DBName
			dbCgf.DBName = ""
			dsn := dbCgf.FormatDSN()
			db = getDBConnection("mysqlWithHooks", dsn)
			err = db.Ping()
			if err != nil {
				log.Panic(err)
			}

			createDBSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS  `%s`;", dbName)
			_, err = db.Exec(createDBSql)
			if err != nil {
				log.Panic(err)
			}

			initSqlStr := fmt.Sprintf("USE  `%s`;", dbName)
			sqlBytes, err := model.SqlFS.ReadFile("mysql.sql")
			if err != nil {
				log.Panic(err)
			}
			initSqlStr = initSqlStr + string(sqlBytes)

			// 分割SQL语句
			sqlStatements := strings.Split(initSqlStr, ";")
			for _, stmt := range sqlStatements {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}
				_, err = db.Exec(stmt)
				if err != nil {
					log.Panicf("执行SQL失败: %s\n错误: %v", stmt, err)
				}
			}
		} else {
			log.Panic(err)
		}
		db := getDBConnection(driverName, conf.Mysql.DSN)
		err := db.Ping()
		if err != nil {
			log.Panic(err)
		}
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	DB = db
}

func getDBConnection(driverName, dataSourceName string) *sqlx.DB {
	sqlDB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	db := sqlx.NewDb(sqlDB, driverName)
	return db
}
