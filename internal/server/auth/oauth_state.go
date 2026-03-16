package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// SignOAuthState creates an HMAC-signed cookie value embedding nonce, site_id, and intent.
// The nonce is sent as the OAuth state parameter; the signed cookie prevents tampering with site_id/intent.
func SignOAuthState(nonce, siteID, intent, secret string) string {
	payload := nonce + "|" + siteID + "|" + intent
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload + "|" + sig
}

// VerifyOAuthState verifies the HMAC signature and extracts site_id and intent.
// expectedNonce is the state parameter returned by the OAuth provider.
func VerifyOAuthState(signed, expectedNonce, secret string) (siteID, intent string, err error) {
	parts := strings.SplitN(signed, "|", 4)
	if len(parts) != 4 {
		return "", "", fmt.Errorf("malformed oauth state")
	}
	nonce, siteID, intent, sig := parts[0], parts[1], parts[2], parts[3]

	if nonce != expectedNonce {
		return "", "", fmt.Errorf("nonce mismatch")
	}

	payload := nonce + "|" + siteID + "|" + intent
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", "", fmt.Errorf("invalid signature")
	}

	return siteID, intent, nil
}
