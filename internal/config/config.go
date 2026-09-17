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
}

type AppConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
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
