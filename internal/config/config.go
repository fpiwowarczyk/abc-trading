package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Port is the port on which the server will listen.
	Port int
	// MaxK is the number of buckets to store in the symbol each bucket is 10 times bigger than the previous one.
	MaxK int
	// MaxBatch is the maximum number of points that can be added in a single batch.
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

// EnvOrDefault returns the value of the environment variable or the default value if the variable is not set.
func EnvOrDefault(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return v
}
