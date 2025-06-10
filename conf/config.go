package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
)

var (
	Conf      *Config
	AppName   string
	RootPath  string
	Languages []string
	Server    ServerConf
	Logger    LoggerConf
	Redis     RedisConf
	Mysql     MysqlConf
	Sqlite    SqliteConf
)

type Config struct {
	AppName   string     `toml:"appName"`
	RootPath  string     `toml:"rootPath"`
	Languages []string   `toml:"languages"`
	Server    ServerConf `toml:"server"`
	Logger    LoggerConf `toml:"logger"`
	Redis     RedisConf  `toml:"redis"`
	Mysql     MysqlConf  `toml:"mysql"`
	Sqlite    SqliteConf `toml:"sqlite"`
}

type ServerConf struct {
	Mode          string `toml:"mode"`
	Address       string `toml:"address"`
	Port          string `toml:"port"`
	EnableMigrate bool   `toml:"enableMigrate"`
	Secret        string `toml:"secret"`
	DBType        string `toml:"dbType"`
}

type LoggerConf struct {
	Level           zapcore.Level `toml:"level"`
	StackTraceLevel zapcore.Level `toml:"stackTraceLevel"`
	Filename        string        `toml:"filename"`
	MaxSize         int           `toml:"maxSize"`
	MaxBackups      int           `toml:"maxBackups"`
	MaxAge          int           `toml:"maxAge"`
	EnableCompress  bool          `toml:"enableCompress"`
}

type MysqlConf struct {
	DSN string `toml:"dsn"`
}

type SqliteConf struct {
	Path string `toml:"path"`
}

type RedisConf struct {
	Enable bool   `toml:"Enable"`
	DSN    string `toml:"dsn"`
	Expire int    `toml:"expire"`
}

func Initialize() {
	conf := NewDefaultConfig()
	viper := viper.New()
	//不设置时，Viper会自动寻找AddConfigPath下的 config.* 配置文件
	//viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		updateGlobalVars()
	})
	viper.WatchConfig()
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
	Conf = conf
	updateGlobalVars()
}

func updateGlobalVars() {
	AppName = Conf.AppName
	RootPath = Conf.RootPath
	Languages = Conf.Languages
	Server = Conf.Server
	Logger = Conf.Logger
	Redis = Conf.Redis
	Mysql = Conf.Mysql
	Sqlite = Conf.Sqlite
}

func NewDefaultConfig() *Config {
	_, filename, _, _ := runtime.Caller(0)
	rootPath := path.Dir(path.Dir(filename))
	server := ServerConf{
		Mode:          "debug",
		Address:       "0.0.0.0",
		Port:          "8888",
		EnableMigrate: false,
		Secret:        "gin-template",
		DBType:        "sqlite",
	}
	languages := []string{
		"en",
		"zh",
	}
	logger := LoggerConf{
		Level:           zapcore.InfoLevel,
		StackTraceLevel: zapcore.ErrorLevel,
		Filename:        "app.log",
		MaxSize:         200,
		MaxBackups:      20,
		MaxAge:          1,
	}

	redis := RedisConf{
		Enable: false,
		DSN:    "rediss://:@localhost:6379",
		Expire: 3600,
	}
	mysql := MysqlConf{
		DSN: "user:pass@tcp(127.0.0.1:3306)/templatedb?charset=utf8mb4",
	}

	sqlite := SqliteConf{
		Path: "app.db",
	}

	return &Config{
		AppName:   "Gin APP Template",
		RootPath:  rootPath,
		Languages: languages,
		Server:    server,
		Logger:    logger,
		Redis:     redis,
		Mysql:     mysql,
		Sqlite:    sqlite,
	}
}
