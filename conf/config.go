package conf

import (
	"embed"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
)

var (
	AppName       string
	RootPath      string
	Timezone      string
	Languages     []string
	Goroutines    int
	Conf          *Config
	Server        ServerConf
	Logger        LoggerConf
	Scheduler     SchedulerConf
	Redis         RedisConf
	DB            DBConf
	Proxy         ProxyConf
	ViperInstance *viper.Viper
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

//go:embed default.toml
var defaultConfigFS embed.FS

func Initialize() {
	ViperInstance = viper.New()
	ViperInstance.SetConfigType("toml")

	// 打开嵌入的配置文件
	defaultConfigFile, err := defaultConfigFS.Open("default.toml")
	if err != nil {
		panic(fmt.Sprintf("Fatal error: failed to open embedded config file: %v", err))
	}
	defer defaultConfigFile.Close()
	err = ViperInstance.ReadConfig(defaultConfigFile)
	if err != nil {
		panic(fmt.Sprintf("Fatal error: failed to read embedded config: %v", err))
	}

	//不设置时，Viper会自动寻找AddConfigPath下的 config.* 配置文件
	rootPath := GetRootPath()
	ViperInstance.SetConfigName("config")
	ViperInstance.AddConfigPath(rootPath)
	ViperInstance.AddConfigPath("../")
	// 查找并合并用户配置文件。如果文件不存在，Viper不会报错，这正是我们想要的。
	// 它会覆盖掉之前从嵌入文件加载的默认值。
	err = ViperInstance.MergeInConfig()
	if err != nil {
		// 如果错误不是 "配置文件未找到"，则打印错误
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			fmt.Println("Warning: failed to merge user config file:", err)
		}
	}

	// 加载环境变量配置
	ViperInstance.AutomaticEnv()
	ViperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var conf Config
	conf.RootPath = rootPath
	if err := ViperInstance.Unmarshal(&conf); err != nil {
		panic(fmt.Sprintf("Fatal error: unable to unmarshal config: %v", err))
	}

	ViperInstance.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := ViperInstance.Unmarshal(&conf); err != nil {
			fmt.Println("Error unmarshalling config after change:", err)
		} else {
			Conf = &conf
			updateGlobalVars()
		}
	})
	ViperInstance.WatchConfig()

	Conf = &conf
	updateGlobalVars()
}

func updateGlobalVars() {
	fmt.Printf("Load Config file: %s\n", ViperInstance.ConfigFileUsed())
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

// GetRootPath 通过探测 go.mod 文件来智能确定项目根目录
func GetRootPath() string {
	// 尝试从当前工作目录向上查找 go.mod
	dir, err := os.Getwd()
	if err != nil {
		return getExecutableDir()
	}

	// 无限循环，向上查找
	for {
		// 检查当前目录下是否存在 go.mod
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		// 到达文件系统根目录，仍未找到
		if dir == filepath.Dir(dir) {
			break
		}

		dir = filepath.Dir(dir)
	}

	return getExecutableDir()
}

// getExecutableDir 获取可执行文件所在的目录
func getExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("无法获取可执行文件路径: %v", err))
	}
	return filepath.Dir(exe)
}
