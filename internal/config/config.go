package config

import "os"

type Config struct {
	Port    string
	Version string
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	ver := os.Getenv("APP_VERSION")
	if ver == "" {
		ver = "dev"
	}
	return Config{Port: port, Version: ver}
}
