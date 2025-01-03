package logging

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

type Mode string

const (
	Development Mode = "development"
	Production  Mode = "production"
)

var ErrUnknownMode = errors.New("unknown mode")

func Load(mode Mode) (*zap.Logger, error) {
	switch mode {
	case Development:
		return zap.NewDevelopment()
	case Production:
		return zap.NewProduction()
	default:
		return nil, ErrUnknownMode
	}
}

func ModeFromEnv(varName string) Mode {
	// TODO: strconv.ParseBool を使う
	if os.Getenv(varName) != "" {
		return Development
	}
	return Production
}
