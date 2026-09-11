package redis

import "github.com/kelseyhightower/envconfig"

// Config holds the Redis connection configuration.
type Config struct {
	Addr     string `default:"localhost:6379" envconfig:"ADDR"`
	Password string `default:"" envconfig:"PASSWORD"`
	DB       int    `default:"0" envconfig:"DB"`
}

// NewConfig loads Redis config from environment variables prefixed with REDIS_.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("redis", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
