package testutil

import (
	"context"
	"testing"

	"github.com/traPtitech/Jomon/internal/logging"
)

func NewContext(t *testing.T) context.Context {
	t.Helper()
	logger := LoadLogger(t)
	ctx := context.Background()
	return logging.SetLogger(ctx, logger)
}
