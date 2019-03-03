package config

import (
	goConfig "github.com/micro/go-config"
)

// ApplicationConfig stores all configuration values
type ApplicationConfig struct {
	Server   Server   `json:"server,omitempty"`
	Database Database `json:"database,omitempty"`
}

type Server struct {
	Host      string `json:"host,omitempty"`
	Port      int    `json:"port,omitempty"`
	StaticDir string `json:"static_dir,omitempty"`
}

type Database struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}

// LoadConfig takes a unknown amount of configs and loads them into the
// ApplicationConfig struct. It returns a pointer to the config struct
func LoadConfig(filePaths ...string) *ApplicationConfig {
	for _, filePath := range filePaths {
		goConfig.LoadFile(filePath)
	}

	cfg := ApplicationConfig{}
	goConfig.Scan(&cfg)

	return &cfg
}
