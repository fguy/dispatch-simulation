package main

import (
	"github.com/fguy/dispatch-simulation/config"

	"go.uber.org/zap"
)

// NewLogger -
func NewLogger(cfg *config.AppConfig) (*zap.Logger, error) {
	if cfg.Logging.Development {
		cfg.Logging.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		cfg.Logging.EncoderConfig = zap.NewProductionEncoderConfig()
	}
	return cfg.Logging.Build()
}
