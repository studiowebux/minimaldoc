package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// SignOAuthState creates an HMAC-signed cookie value embedding nonce, site_id, intent, and PKCE verifier.
// The nonce is sent as the OAuth state parameter; the signed cookie prevents tampering.
func SignOAuthState(nonce, siteID, intent, secret string) string {
	return signOAuthStateWithVerifier(nonce, siteID, intent, "", secret)
}

// SignOAuthStateWithPKCE creates an HMAC-signed cookie value that also includes a PKCE code verifier.
func SignOAuthStateWithPKCE(nonce, siteID, intent, verifier, secret string) string {
	return signOAuthStateWithVerifier(nonce, siteID, intent, verifier, secret)
}

func signOAuthStateWithVerifier(nonce, siteID, intent, verifier, secret string) string {
	payload := nonce + "|" + siteID + "|" + intent + "|" + verifier
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload + "|" + sig
}

// VerifyOAuthState verifies the HMAC signature and extracts site_id and intent.
// expectedNonce is the state parameter returned by the OAuth provider.
func VerifyOAuthState(signed, expectedNonce, secret string) (siteID, intent string, err error) {
	_, siteID, intent, err = VerifyOAuthStateWithPKCE(signed, expectedNonce, secret)
	return
}

// VerifyOAuthStateWithPKCE verifies the HMAC signature and returns site_id, intent, and PKCE verifier.
func VerifyOAuthStateWithPKCE(signed, expectedNonce, secret string) (verifier, siteID, intent string, err error) {
	parts := strings.SplitN(signed, "|", 5)
	if len(parts) < 4 {
		return "", "", "", fmt.Errorf("malformed oauth state")
	}

	// Support both old 4-part format and new 5-part format with verifier
	var nonce, sig string
	verifier = ""
	if len(parts) == 5 {
		nonce, siteID, intent, verifier, sig = parts[0], parts[1], parts[2], parts[3], parts[4]
	} else {
		nonce, siteID, intent, sig = parts[0], parts[1], parts[2], parts[3]
	}

	if nonce != expectedNonce {
		return "", "", "", fmt.Errorf("nonce mismatch")
	}

	payload := nonce + "|" + siteID + "|" + intent + "|" + verifier
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", "", "", fmt.Errorf("invalid signature")
	}

	return verifier, siteID, intent, nil
}

// GeneratePKCEVerifier creates a random PKCE code verifier (43-128 chars, RFC 7636).
func GeneratePKCEVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate PKCE verifier: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// PKCECodeChallenge computes the S256 code challenge for a given verifier (RFC 7636).
func PKCECodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
