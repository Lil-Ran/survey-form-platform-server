package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig holds the configuration values from application.yml
type AppConfig struct {
	Datasource struct {
		DriverName string `mapstructure:"driverName"`
		Host       string `mapstructure:"host"`
		Port       int    `mapstructure:"port"`
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		DBName     string `mapstructure:"dbname"`
		Charset    string `mapstructure:"charset"`
	} `mapstructure:"datasource"`

	Server struct {
		Host        string `mapstructure:"host"`
		Port        int    `mapstructure:"port"`
		EnableHTTPS bool   `mapstructure:"enable_https"`
		LogLevel    string `mapstructure:"log_level"`
	} `mapstructure:"server"`

	Auth struct {
		JWTSecret   string `mapstructure:"jwt_secret"`
		TokenExpiry string `mapstructure:"token_expiry"`
	} `mapstructure:"auth"`

	SMTP struct {
		From     string `mapstructure:"from"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"smtp"`
}

// Global variable to hold the loaded configuration
var Config AppConfig

// LoadConfig loads the configuration from application.yml
func LoadConfig() {
	viper.SetConfigName("application") // 配置文件名称 (不带扩展名)
	viper.SetConfigType("yml")         // 配置文件类型
	viper.AddConfigPath("./config")    // 配置文件路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// 解析配置文件
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Error unmarshaling config, %s", err)
	}
}
