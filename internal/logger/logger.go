package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	ServiceName string
	InstanceID  string
	Output      io.Writer
}

var appLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func Init(cfg Config) {
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}

	context := zerolog.New(output).With().Timestamp()
	if cfg.ServiceName != "" {
		context = context.Str("service", cfg.ServiceName)
	}
	if cfg.InstanceID != "" {
		context = context.Str("instance_id", cfg.InstanceID)
	}

	appLogger = context.Logger()
}

func L() *zerolog.Logger {
	return &appLogger
}
