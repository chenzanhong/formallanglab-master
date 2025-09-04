package configs

import (
	"os"
	"gopkg.in/yaml.v3"
)

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	PG PGConfig `yaml:"pg"`
}

func LoadConfig() (*Config, error) {
	yamlFile, err := os.ReadFile("./config.yaml")
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
