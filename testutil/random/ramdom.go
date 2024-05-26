package random

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// AlphaNumeric 指定した文字数のランダム英数字文字列を生成します
// この関数はmath/randが生成する擬似乱数を使用します
func AlphaNumeric(t *testing.T, n int) string {
	t.Helper()
	b := make([]byte, n)
	_, err := rand.Read(b)
	require.NoError(t, err)

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}

	return result
}

func Numeric(t *testing.T, max int) int {
	t.Helper()
	n := rand.Intn(max)
	return n
}

func Numeric64(t *testing.T, max int64) int64 {
	t.Helper()
	n := rand.Int63n(max)
	return n
}

func AlphaNumericSlice(t *testing.T, length int, max int64) []string {
	slice := []string{}
	for range length {
		slice = append(slice, AlphaNumeric(t, int(max)))
	}
	return slice
}

func NumericSlice(t *testing.T, length int, max int) []int {
	slice := []int{}
	for range length {
		slice = append(slice, Numeric(t, max))
	}
	return slice
}

func Numeric64Slice(t *testing.T, length int, max int64) []int64 {
	slice := []int64{}
	for range length {
		slice = append(slice, Numeric64(t, max))
	}
	return slice
}
