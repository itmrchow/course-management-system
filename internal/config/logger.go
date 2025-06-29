package config

import (
	"os"

	"github.com/rs/zerolog"
)

func InitLogger() *zerolog.Logger {
	// TODO: setting log
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger
}
