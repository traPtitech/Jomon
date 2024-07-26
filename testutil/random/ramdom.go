package random

import (
	"math/rand/v2"
	"testing"

	"github.com/samber/lo"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// AlphaNumeric 指定した文字数のランダム英数字文字列を生成します
// この関数はmath/randが生成する擬似乱数を使用します
func AlphaNumeric(t *testing.T, n int) string {
	t.Helper()
	b := make([]byte, n)

	for i := range n {
		b[i] = letters[rand.IntN(len(letters))]
	}

	return string(b)
}

func Numeric(t *testing.T, max int) int {
	t.Helper()
	return rand.IntN(max)
}

func Numeric64(t *testing.T, max int64) int64 {
	t.Helper()
	return rand.Int64N(max)
}

func AlphaNumericSlice(t *testing.T, length int, max int64) []string {
	return lo.Times(length, func(index int) string {
		return AlphaNumeric(t, int(max))
	})
}

func NumericSlice(t *testing.T, length int, max int) []int {
	return lo.Times(length, func(index int) int {
		return Numeric(t, max)
	})
}

func Numeric64Slice(t *testing.T, length int, max int64) []int64 {
	return lo.Times(length, func(index int) int64 {
		return Numeric64(t, max)
	})
}
