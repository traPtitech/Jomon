package random

import (
	"math/rand/v2"
	"testing"
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
	t.Helper()
	slice := []string{}
	for range length {
		slice = append(slice, AlphaNumeric(t, int(max)))
	}
	return slice
}

func NumericSlice(t *testing.T, length int, max int) []int {
	t.Helper()
	slice := []int{}
	for range length {
		slice = append(slice, Numeric(t, max))
	}
	return slice
}

func Numeric64Slice(t *testing.T, length int, max int64) []int64 {
	t.Helper()
	slice := []int64{}
	for range length {
		slice = append(slice, Numeric64(t, max))
	}
	return slice
}
