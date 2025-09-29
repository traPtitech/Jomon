package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/internal/logging"
	"go.uber.org/zap"
)

func LoadLogger(t *testing.T) *zap.Logger {
	t.Helper()
	logger, err := logging.Load(logging.Development)
	require.NoError(t, err)
	return logger
}
