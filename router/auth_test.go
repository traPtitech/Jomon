package router

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sample verifier and challenge comes from Appendix B of RFC7636
// See: https://tools.ietf.org/html/rfc7636#appendix-B
func TestHandler_GetCodeChallenge(t *testing.T) {
	codeVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	codeVerifierHash := sha256.Sum256([]byte(codeVerifier))

	encoder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_").WithPadding(base64.NoPadding)
	challenge := encoder.EncodeToString(codeVerifierHash[:])

	assert.Equal(t, "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM", challenge)
}
