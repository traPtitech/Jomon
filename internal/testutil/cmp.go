package testutil

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func ApproxEqualOptions() []cmp.Option {
	return []cmp.Option{
		cmpopts.EquateApproxTime(time.Second),
		cmpopts.EquateEmpty(),
	}
}

func AssertEqual(t *testing.T, expected, actual interface{}, opts ...cmp.Option) bool {
	t.Helper()
	diff := cmp.Diff(expected, actual, opts...)
	if diff == "" {
		return true
	}
	t.Errorf("Not equal (-expected +actual):\n"+
		"expected: %s\n"+
		"actual  : %s\n%s", expected, actual, diff)
	return false
}

func RequireEqual(t *testing.T, expected, actual interface{}, opts ...cmp.Option) {
	t.Helper()
	if !AssertEqual(t, expected, actual, opts...) {
		t.FailNow()
	}
}
