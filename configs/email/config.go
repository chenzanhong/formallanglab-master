package email

import (
	"backend/configs"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"go.yaml.in/yaml/v2"
)

type EmailWorkerConfig struct {
	Email      configs.EMAILConfig      `yaml:"email"`
	SMTPServer configs.SMTPServerConfig `yaml:"smtp_server"`
	Kafka      configs.KafkaConfig      `yaml:"kafka"`
}

// LoadConfig 加载配置文件并返回 DBConfig
func LoadEmailWorkerConfig() (*EmailWorkerConfig, error) {
	_, filename, _, ok := runtime.Caller(0) // 获取当前的文件名
	if !ok {
		log.Fatal("无法获取运行时调用者信息")
	}

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(filename)

	// 构建到项目根目录的相对路径
	configPath := filepath.Join(currentDir, "config.yaml")

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config EmailWorkerConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SetEmailEnvVariables() {
	config, err := LoadEmailWorkerConfig()
	if err != nil {
		log.Fatalf("加载配置失败：%v", err.Error())
	}

	// 辅助函数：如果 envVar 未设置，则用 fallback 值设置它
	setEnvIfNotSet := func(envVar, fallback string) {
		if os.Getenv(envVar) == "" {
			os.Setenv(envVar, fallback)
		}
	}

	// Email
	setEnvIfNotSet("EMAIL_NAME", config.Email.Name)
	setEnvIfNotSet("EMAIL_PASSWORD", config.Email.Password)

	// SMTP Server
	setEnvIfNotSet("SMTP_SERVER_HOST", config.SMTPServer.Host)
	setEnvIfNotSet("SMTP_SERVER_PORT", config.SMTPServer.Port)

	// Kafka
	setEnvIfNotSet("KAFKA_BROKER", config.Kafka.Brokers)
	setEnvIfNotSet("KAFKA_TOPIC", config.Kafka.Topic)
}
