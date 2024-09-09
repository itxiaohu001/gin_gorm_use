package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// AppConfig 全局配置变量
var AppConfig Config

// Config 应用程序配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Storage  StorageConfig  `yaml:"storage"`
	Logger   LoggerConfig   `yaml:"logger"`
	Redis    RedisConfig    `yaml:"redis"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type      string `yaml:"type"`
	LocalPath string `yaml:"local_path"`
	Minio     MinioConfig `yaml:"minio"`
}

// MinioConfig Minio配置
type MinioConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	UseSSL    bool   `yaml:"use_ssl"`
	Bucket    string `yaml:"bucket"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// LoadConfig 加载配置文件
func LoadConfig() error {
	configFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	err = yaml.Unmarshal(configFile, &AppConfig)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 添加调试日志
	fmt.Printf("加载的配置: %+v\n", AppConfig)
	fmt.Printf("存储配置: %+v\n", AppConfig.Storage)
	fmt.Printf("Minio 配置: %+v\n", AppConfig.Storage.Minio)

	return nil
}
