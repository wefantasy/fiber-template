package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strings"
)

var (
	AppName    string
	RootPath   string
	Timezone   string
	Languages  []string
	Goroutines int
	Conf       *Config
	Server     ServerConf
	Logger     LoggerConf
	Scheduler  SchedulerConf
	Redis      RedisConf
	DB         DBConf
	Proxy      ProxyConf
)

type Config struct {
	AppName    string        `toml:"appName"`    // 应用名称
	RootPath   string        `toml:"rootPath"`   // 应用根目录
	Timezone   string        `toml:"timezone"`   // 时区
	Languages  []string      `toml:"languages"`  // 支持的语言
	Goroutines int           `toml:"goroutines"` // 默认协程数量
	Server     ServerConf    `toml:"server"`
	Logger     LoggerConf    `toml:"logger"`
	Scheduler  SchedulerConf `toml:"scheduler"`
	DB         DBConf        `toml:"db"`
	Redis      RedisConf     `toml:"redis"`
	Proxy      ProxyConf     `toml:"proxy"`
}

type ServerConf struct {
	Address string `toml:"address"` // 监听地址
	Port    string `toml:"port"`    // 监听端口
	Secret  string `toml:"secret"`  // jwt密钥/Secret模式密钥
}

type LoggerConf struct {
	Level           zapcore.Level `toml:"level"`           // 日志级别，支持debug(-1)/info(0)/warn(1)/error(2)/dpanic(3)/panic(4)/fatal(5)
	StackTraceLevel zapcore.Level `toml:"stackTraceLevel"` // 堆栈级别
	Filename        string        `toml:"filename"`        // 日志文件路径
	MaxSize         int           `toml:"maxSize"`         // 日志文件最大大小（MB）
	MaxBackups      int           `toml:"maxBackups"`      // 日志文件最大备份数
	MaxAge          int           `toml:"maxAge"`          // 切割后的日志文件最大保留天数
	EnableCompress  bool          `toml:"enableCompress"`  // 是否启用压缩
}

type SchedulerConf struct {
	EnableTasks       []string `toml:"enableTasks"`       // 启用的任务
	RunAtStartupTasks []string `toml:"runAtStartupTasks"` // 启动时运行的任务
}

type DBConf struct {
	Type          string `toml:"type"`          // 数据库类型：mysql, sqlite
	DSN           string `toml:"dsn"`           // 数据库连接字符串: user:pass@tcp(127.0.0.1:3306)/template-db?charset=utf8mb4&parseTime=true, file:test.db?cache=shared&mode=memory
	EnableMigrate bool   `toml:"enableMigrate"` // 是否启用数据库迁移
}

type RedisConf struct {
	Enable bool   `toml:"enable"` // 是否启用
	DSN    string `toml:"dsn"`    // Redis 连接字符串
	Expire int    `toml:"expire"` // 过期时间（秒）
}

type ProxyConf struct {
	BaseUrl string `toml:"baseUrl"` // 代理服务地址
	Secret  string `toml:"secret"`  // 代理服务密钥
}

func Initialize() {
	conf := NewDefaultConfig()
	viperInstance := viper.New()

	// 加载默认配置
	var defaultConfigMap map[string]interface{}
	err := mapstructure.Decode(conf, &defaultConfigMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode default config: %v", err))
	}
	err = viperInstance.MergeConfigMap(defaultConfigMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to merge default config map: %v", err))
	}

	// 加载环境变量配置
	viperInstance.AutomaticEnv()
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//不设置时，Viper会自动寻找AddConfigPath下的 config.* 配置文件
	viperInstance.AddConfigPath(".")
	viperInstance.AddConfigPath(conf.RootPath)
	viperInstance.AddConfigPath("../")
	err = viperInstance.ReadInConfig()
	if err == nil {
		viperInstance.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			if err := viperInstance.Unmarshal(&conf); err != nil {
				fmt.Println("Error unmarshalling config after change:", err)
			}
			Conf = conf
			updateGlobalVars()
		})
		viperInstance.WatchConfig()
	} else {
		fmt.Println(err)
	}
	err = viperInstance.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
	Conf = conf
	updateGlobalVars()
}

func updateGlobalVars() {
	AppName = Conf.AppName
	RootPath = Conf.RootPath
	Timezone = Conf.Timezone
	Languages = Conf.Languages
	Goroutines = Conf.Goroutines
	Server = Conf.Server
	Logger = Conf.Logger
	Scheduler = Conf.Scheduler
	DB = Conf.DB
	Redis = Conf.Redis
	Proxy = Conf.Proxy
}

func NewDefaultConfig() *Config {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootPath := path.Dir(exePath)

	server := ServerConf{
		Address: "0.0.0.0",
		Port:    "8888",
		Secret:  "FiberTemplate",
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

	s := SchedulerConf{
		EnableTasks:       []string{},
		RunAtStartupTasks: []string{},
	}

	db := DBConf{
		Type:          "",
		DSN:           "",
		EnableMigrate: false,
	}

	redis := RedisConf{
		Enable: false,
		DSN:    "rediss://:@localhost:6379",
		Expire: 3600,
	}

	proxy := ProxyConf{
		BaseUrl: "",
		Secret:  "",
	}

	return &Config{
		AppName:  "Fiber APP Template",
		RootPath: rootPath,
		Timezone: "Asia/Shanghai",
		Languages: []string{
			"en",
			"zh",
		},
		Goroutines: 50,
		Server:     server,
		Logger:     logger,
		Scheduler:  s,
		DB:         db,
		Redis:      redis,
		Proxy:      proxy,
	}
}
