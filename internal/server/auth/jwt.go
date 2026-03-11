package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT token claims.
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	SiteID string `json:"sid"`
}

// GenerateToken creates a new JWT token.
func GenerateToken(userID, email, role, siteID, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "minimaldoc-server",
		},
		UserID: userID,
		Email:  email,
		Role:   role,
		SiteID: siteID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken parses and validates a JWT token.
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GenerateRefreshToken creates a longer-lived refresh token.
func GenerateRefreshToken(userID, secret string, expiry time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "minimaldoc-server",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateRefreshToken parses and validates a refresh token.
func ValidateRefreshToken(tokenString, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	return claims.Subject, nil
}
