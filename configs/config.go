package configs

import (
	"backend/logs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"gopkg.in/yaml.v3"
)

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RedisConfig struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Password    string `yaml:"password"`
	DB          int    `yaml:"db"`
	MaxConn     int    `yaml:"max_conn"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
}

type EMAILConfig struct {
	Name     string `yaml:"email_name"`
	Password string `yaml:"email_password"`
}
type SMTPServerConfig struct {
	Host string `yaml:"SMTPServer_host"`
	Port string `yaml:"SMTPServer_port"`
}

type RateConfig struct {
	UserRate  int `yaml:"user_rate"`
	UserBurst int `yaml:"user_burst"`
}

type Config struct {
	PG         PGConfig         `yaml:"pg"`
	Redis      RedisConfig      `yaml:"redis"`
	Email      EMAILConfig      `yaml:"email"`
	SMTPServer SMTPServerConfig `yaml:"smtp_server"`
	Rate       RateConfig       `yaml:"rate"`
}

// getDBConfigPath 获取数据库配置文件的路径
func getDBConfigPath() string {
	_, filename, _, ok := runtime.Caller(2) // 获取调用者的文件名
	if !ok {
		log.Fatal("无法获取运行时调用者信息")
	}

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(filename)

	// 构建到项目根目录的相对路径
	dbConfigPath := filepath.Join(currentDir, "..", "configs", "config.yaml")

	// 将路径转换为绝对路径并简化路径
	absPath, err := filepath.Abs(dbConfigPath)
	if err != nil {
		log.Printf("无法获取绝对路径: %v", err)
	}

	simplifiedPath := filepath.Clean(absPath)

	return simplifiedPath
}

// GetDBConfigPath 返回数据库配置文件的路径
func GetDBConfigPath() string {
	return getDBConfigPath()
}

// LoadConfig 加载配置文件并返回 DBConfig
func LoadConfig() (*Config, error) {
	configPath := GetDBConfigPath()
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SetEnvVariables() {
	config, err := LoadConfig()
	if err != nil {
		logs.Sugar.Fatalf("加载配置失败：%v", err.Error())
	}

	// 辅助函数：如果 envVar 未设置，则用 fallback 值设置它
	setEnvIfNotSet := func(envVar, fallback string) {
		if os.Getenv(envVar) == "" {
			os.Setenv(envVar, fallback)
		}
	}

	// PostgreSQL
	setEnvIfNotSet("DB_USER", config.PG.User)
	setEnvIfNotSet("DB_PASSWORD", config.PG.Password)
	setEnvIfNotSet("DB_HOST", config.PG.Host)
	setEnvIfNotSet("DB_PORT", config.PG.Port)
	setEnvIfNotSet("DB_NAME", config.PG.Name)

	// Redis
	setEnvIfNotSet("REDIS_HOST", config.Redis.Host)
	setEnvIfNotSet("REDIS_PORT", config.Redis.Port)
	setEnvIfNotSet("REDIS_PASSWORD", config.Redis.Password)
	setEnvIfNotSet("REDIS_DB", strconv.Itoa(config.Redis.DB))
	setEnvIfNotSet("REDIS_MAX_CONN", strconv.Itoa(config.Redis.MaxConn))
	setEnvIfNotSet("REDIS_MAX_IDLE_CONN", strconv.Itoa(config.Redis.MaxIdleConn))

	// Email
	setEnvIfNotSet("EMAIL_NAME", config.Email.Name)
	setEnvIfNotSet("EMAIL_PASSWORD", config.Email.Password)

	// SMTP Server
	setEnvIfNotSet("SMTP_SERVER_HOST", config.SMTPServer.Host)
	setEnvIfNotSet("SMTP_SERVER_PORT", config.SMTPServer.Port)

	// Rate Limiting
	setEnvIfNotSet("RATE_USER_RATE", strconv.Itoa(config.Rate.UserRate))
	setEnvIfNotSet("RATE_USER_BURST", strconv.Itoa(config.Rate.UserBurst))
}
