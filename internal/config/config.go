package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port     int
	MaxK     int
	MaxBatch int
}

func NewConfig() (*Config, error) {
	port, err := strconv.Atoi(EnvOrDefault("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	maxK, err := strconv.Atoi(EnvOrDefault("MAX_K", "8"))
	if err != nil {
		return nil, err
	}

	maxBatch, err := strconv.Atoi(EnvOrDefault("MAX_BATCH", "10000"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:     port,
		MaxK:     maxK,
		MaxBatch: maxBatch,
	}, nil
}

func EnvOrDefault(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return v
}
