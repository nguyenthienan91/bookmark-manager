package api

import ("github.com/kelseyhightower/envconfig"
	"github.com/google/uuid"
)

type Config struct {
	AppPort string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-service" envconfig:"SERVICE_NAME"`
	InstanceID string `envconfig:"INSTANCE_ID"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("api", cfg)
	if err != nil {
		return nil, err
	} 
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.NewString()
	} 	
	return cfg, nil
}