package config

var Cfg Config

type Config struct {
	DatabaseConfig   *DatabaseConfig
	StorageAccessKey string
	StorageSecretKey string
	ChatGPTKey       string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}
