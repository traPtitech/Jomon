package testutil

import (
	"context"
	"testing"
)

func NewContext(t *testing.T) context.Context {
	t.Helper()
	return context.Background()
}
