package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// Context extraction errors
var (
	errMissingSiteID = errors.New("site_id not found in context")
	errMissingUserID = errors.New("user_id not found in context")
	errMissingUser   = errors.New("user not found in context")
	errMissingRole   = errors.New("user_role not found in context")
)

// getSiteID extracts site_id from gin context.
func getSiteID(c *gin.Context) (string, error) {
	val, ok := c.Get("site_id")
	if !ok {
		return "", errMissingSiteID
	}
	siteID, ok := val.(string)
	if !ok {
		return "", errMissingSiteID
	}
	return siteID, nil
}

// getUserID extracts user_id from gin context.
func getUserID(c *gin.Context) (string, error) {
	val, ok := c.Get("user_id")
	if !ok {
		return "", errMissingUserID
	}
	userID, ok := val.(string)
	if !ok {
		return "", errMissingUserID
	}
	return userID, nil
}

// getUserClaims extracts user claims from gin context.
func getUserClaims(c *gin.Context) (*auth.Claims, error) {
	val, ok := c.Get("user")
	if !ok {
		return nil, errMissingUser
	}
	claims, ok := val.(*auth.Claims)
	if !ok {
		return nil, errMissingUser
	}
	return claims, nil
}

// getUserRole extracts user_role from gin context.
func getUserRole(c *gin.Context) (string, error) {
	val, ok := c.Get("user_role")
	if !ok {
		return "", errMissingRole
	}
	role, ok := val.(string)
	if !ok {
		return "", errMissingRole
	}
	return role, nil
}
