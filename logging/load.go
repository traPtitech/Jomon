package logging

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Mode string

const (
	Development Mode = "development"
	Production  Mode = "production"
)

var ErrUnknownMode = errors.New("unknown mode")

func Load(mode Mode) (*zap.Logger, error) {
	var config zap.Config
	switch mode {
	case Development:
		config = zap.NewDevelopmentConfig()
	case Production:
		config = zap.NewProductionConfig()
	default:
		return nil, ErrUnknownMode
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	return config.Build()
}

func ModeFromEnv(varName string) Mode {
	// TODO: strconv.ParseBool を使う
	if os.Getenv(varName) != "" {
		return Development
	}
	return Production
}
