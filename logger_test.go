package main

import (
	"testing"

	"github.com/fguy/dispatch-simulation/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger_Development(t *testing.T) {
	t.Parallel()

	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.EncoderConfig = zapcore.EncoderConfig{}

	_, err := NewLogger(&config.AppConfig{
		Logging: loggerCfg,
	})
	assert.NoError(t, err)
}

func TestNewLogger_Production(t *testing.T) {
	t.Parallel()

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.EncoderConfig = zapcore.EncoderConfig{}

	_, err := NewLogger(&config.AppConfig{
		Logging: loggerCfg,
	})
	assert.NoError(t, err)
}
