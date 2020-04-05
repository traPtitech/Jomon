package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Sample verifier and challenge comes from Appendix B of RFC7636
// See: https://tools.ietf.org/html/rfc7636#appendix-B
func TestGetCodeChallenge(t *testing.T) {
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	challenge := getCodeChallenge([]byte(verifier))

	assert.Equal(t, "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM", challenge)
}
