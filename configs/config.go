package configs

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type EMAILConfig struct {
	Name     string `yaml:"email_name"`
	Password string `yaml:"email_password"`
}
type SMTPServerConfig struct {
	Host string `yaml:"SMTPServer_host"`
	Port string `yaml:"SMTPServer_port"`
}

type Config struct {
	PG         PGConfig         `yaml:"pg"`
	Email      EMAILConfig      `yaml:"email"`
	SMTPServer SMTPServerConfig `yaml:"smtp_server"`
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
