package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	App       AppConfig       `json:"app"`
	Database  DatabaseConfig  `json:"database"`
	JWT       JWTConfig       `json:"jwt"`
	Migration MigrationConfig `json:"migration"`
	Storage   StorageConfig   `json:"storage"`
	CORS      CORSConfig      `json:"cors"`
}

type AppConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	LogLevel string `json:"log_level"`
}

type DatabaseConfig struct {
	Driver    string `json:"driver"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Charset   string `json:"charset"`
	ParseTime bool   `json:"parse_time"`
}

type JWTConfig struct {
	Secret string `json:"secret"`
	Exp    int    `json:"exp"`
}

type MigrationConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type StorageConfig struct {
	Path string `json:"path"`
}

type CORSConfig struct {
	AllowOrigins     []string `json:"allow_origins"`
	AllowMethods     []string `json:"allow_methods"`
	AllowHeaders     []string `json:"allow_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
}


func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
