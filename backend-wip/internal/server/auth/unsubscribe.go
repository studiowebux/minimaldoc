package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// SignUnsubscribeToken creates an HMAC-signed token for newsletter unsubscribe links.
// The token encodes siteID + email so the unsubscribe endpoint can verify the request
// came from a legitimate email recipient, preventing unauthorized unsubscribes.
func SignUnsubscribeToken(siteID, email, secret string) string {
	payload := siteID + "|" + email
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}

// VerifyUnsubscribeToken verifies the HMAC signature and extracts siteID and email.
func VerifyUnsubscribeToken(token, secret string) (siteID, email string, err error) {
	// Split into payload and signature
	dot := -1
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '.' {
			dot = i
			break
		}
	}
	if dot < 0 {
		return "", "", fmt.Errorf("malformed unsubscribe token")
	}

	encodedPayload := token[:dot]
	sig := token[dot+1:]

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return "", "", fmt.Errorf("malformed unsubscribe token")
	}
	payload := string(payloadBytes)

	// Verify HMAC
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payloadBytes)
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", "", fmt.Errorf("invalid unsubscribe token")
	}

	// Extract siteID and email
	pipe := -1
	for i, c := range payload {
		if c == '|' {
			pipe = i
			break
		}
	}
	if pipe < 0 {
		return "", "", fmt.Errorf("malformed unsubscribe token payload")
	}

	return payload[:pipe], payload[pipe+1:], nil
}
