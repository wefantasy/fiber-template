package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"path/filepath"
)

var (
	Conf   *Config
	Base   BaseConf
	Server ServerConf
	Logger LoggerConf
	Redis  RedisConf
	Mysql  MysqlConf
	Sqlite SqliteConf
)

type Config struct {
	Base   BaseConf   `toml:"base"`
	Server ServerConf `toml:"server"`
	Logger LoggerConf `toml:"logger"`
	Redis  RedisConf  `toml:"redis"`
	Mysql  MysqlConf  `toml:"mysql"`
	Sqlite SqliteConf `toml:"sqlite"`
}

type BaseConf struct {
	AppName     string   `toml:"appName"`     // 应用名称
	RootPath    string   `toml:"rootPath"`    // 应用根目录
	Languages   []string `toml:"languages"`   // 支持的语言
	ProxyApi    string   `toml:"proxyApi"`    // 代理服务地址API
	ProxySecret string   `toml:"proxySecret"` // 代理密钥
}

type ServerConf struct {
	Address       string `toml:"address"`       // 监听地址
	Port          string `toml:"port"`          // 监听端口
	EnableMigrate bool   `toml:"enableMigrate"` // 是否启用数据库迁移
	Secret        string `toml:"secret"`        // jwt密钥/Secret模式密钥
	DBType        string `toml:"dbType"`        // 数据库类型：mysql, sqlite
}

type LoggerConf struct {
	Level           zapcore.Level `toml:"level"`           // 日志级别
	StackTraceLevel zapcore.Level `toml:"stackTraceLevel"` // 堆栈级别
	Filename        string        `toml:"filename"`        // 日志文件路径
	MaxSize         int           `toml:"maxSize"`         // 日志文件最大大小（MB）
	MaxBackups      int           `toml:"maxBackups"`      // 日志文件最大备份数
	MaxAge          int           `toml:"maxAge"`          // 切割后的日志文件最大保留天数
	EnableCompress  bool          `toml:"enableCompress"`  // 是否启用压缩
}

type MysqlConf struct {
	DSN string `toml:"dsn"` // 数据库连接字符串
}

type SqliteConf struct {
	Path string `toml:"path"` // 数据库文件路径
}

type RedisConf struct {
	Enable bool   `toml:"enable"` // 是否启用
	DSN    string `toml:"dsn"`    // Redis 连接字符串
	Expire int    `toml:"expire"` // 过期时间（秒）
}

func Initialize() {
	conf := NewDefaultConfig()
	viperInstance := viper.New()
	//不设置时，Viper会自动寻找AddConfigPath下的 config.* 配置文件
	viperInstance.AddConfigPath(".")
	viperInstance.AddConfigPath(conf.Base.RootPath)
	viperInstance.AddConfigPath("../")
	err := viperInstance.ReadInConfig()
	if err == nil {
		viperInstance.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			updateGlobalVars()
		})
		viperInstance.WatchConfig()
		err = viperInstance.Unmarshal(&conf)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println(err)
	}
	Conf = conf
	updateGlobalVars()
}

func updateGlobalVars() {
	fmt.Println("Use config: ")
	fmt.Printf("%+v", Conf)
	Base = Conf.Base
	Server = Conf.Server
	Logger = Conf.Logger
	Redis = Conf.Redis
	Mysql = Conf.Mysql
	Sqlite = Conf.Sqlite
}

func NewDefaultConfig() *Config {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootPath := path.Dir(exePath)

	base := BaseConf{
		AppName:  "Fiber APP Template",
		RootPath: rootPath,
		Languages: []string{
			"en",
			"zh",
		},
		ProxyApi:    "",
		ProxySecret: "",
	}

	server := ServerConf{
		Address:       "0.0.0.0",
		Port:          "8888",
		EnableMigrate: false,
		Secret:        "FiberTemplate",
		DBType:        "",
	}

	logger := LoggerConf{
		Level:           zapcore.InfoLevel,
		StackTraceLevel: zapcore.ErrorLevel,
		Filename:        "app.log",
		MaxSize:         100,
		MaxBackups:      30,
		MaxAge:          30,
		EnableCompress:  true,
	}

	redis := RedisConf{
		Enable: false,
		DSN:    "rediss://:@localhost:6379",
		Expire: 3600,
	}
	mysql := MysqlConf{
		DSN: "user:pass@tcp(127.0.0.1:3306)/templateDb?charset=utf8mb4&parseTime=true",
	}

	sqlite := SqliteConf{
		Path: filepath.Join(rootPath, "app.db"),
	}

	return &Config{
		Base:   base,
		Server: server,
		Logger: logger,
		Redis:  redis,
		Mysql:  mysql,
		Sqlite: sqlite,
	}
}
