package config

import (
	"fmt"
	"strings"
)

// Config is the server configuration read from the environment.
type Config struct {
	Port     string
	Debug    bool
	CertPath string
	KeyPath  string
}

// Parse reads HTTP_PORT, DEBUG, and TLS paths. HTTP_PORT is required.
// TLS paths are required unless DEBUG is exactly "true".
func Parse(getenv func(string) string) (Config, error) {
	cfg := Config{
		Port:     strings.TrimSpace(getenv("HTTP_PORT")),
		Debug:    getenv("DEBUG") == "true",
		CertPath: strings.TrimSpace(getenv("SSL_CERT_PATH")),
		KeyPath:  strings.TrimSpace(getenv("SSL_KEY_PATH")),
	}
	if cfg.Port == "" {
		return Config{}, fmt.Errorf("HTTP_PORT is required")
	}
	if !cfg.Debug {
		if cfg.CertPath == "" {
			return Config{}, fmt.Errorf("SSL_CERT_PATH is required when DEBUG is not true")
		}
		if cfg.KeyPath == "" {
			return Config{}, fmt.Errorf("SSL_KEY_PATH is required when DEBUG is not true")
		}
	}
	return cfg, nil
}
