package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error codes follow UPPER_SNAKE_CASE convention.
// Each error has three facets:
//   1. code  — stable identifier for programmatic use
//   2. message — friendly message for end users
//   3. HTTP status — appropriate status code

// General
const (
	ErrBadRequest     = "BAD_REQUEST"
	ErrDatabaseError  = "DATABASE_ERROR"
	ErrInternalError  = "INTERNAL_ERROR"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrForbidden      = "FORBIDDEN"
	ErrNotFound       = "NOT_FOUND"
	ErrConflict       = "CONFLICT"
	ErrRateLimited    = "RATE_LIMITED"
	ErrNotAvailable   = "NOT_AVAILABLE"
	ErrMissingParams  = "MISSING_PARAMETERS"
	ErrInvalidRequest = "INVALID_REQUEST"
)

// Authentication & Authorization
const (
	ErrInvalidCredentials     = "INVALID_CREDENTIALS"
	ErrInvalidBootstrapToken  = "INVALID_BOOTSTRAP_TOKEN"
	ErrInvalidRefreshToken    = "INVALID_REFRESH_TOKEN"
	ErrInvalidState           = "INVALID_STATE"
	ErrInvalidVerifyToken     = "INVALID_VERIFY_TOKEN"
	ErrTokenGenerationFailed  = "TOKEN_GENERATION_FAILED"
	ErrPasswordRequired       = "PASSWORD_REQUIRED"
	ErrPasswordTooShort       = "PASSWORD_TOO_SHORT"
	ErrPasswordHashFailed     = "PASSWORD_HASH_FAILED"
	ErrUserNotFound           = "USER_NOT_FOUND"
	ErrUserAlreadyExists      = "USER_ALREADY_EXISTS"
	ErrUserCreationFailed     = "USER_CREATION_FAILED"
	ErrEmailVerifyFailed      = "EMAIL_VERIFY_FAILED"
	ErrRegistrationRequired   = "REGISTRATION_REQUIRED"
	ErrUnknownProvider        = "UNKNOWN_PROVIDER"
	ErrOAuthLinkFailed        = "OAUTH_LINK_FAILED"
	ErrMissingSiteContext     = "MISSING_SITE_CONTEXT"
	ErrCannotModifySelf       = "CANNOT_MODIFY_SELF"
	ErrInsufficientPermission = "INSUFFICIENT_PERMISSION"
	ErrOwnPostsOnly           = "OWN_POSTS_ONLY"
	ErrBanned                 = "BANNED"
)

// Site
const (
	ErrSiteNotFound       = "SITE_NOT_FOUND"
	ErrSiteInvalid        = "SITE_INVALID"
	ErrSiteCreationFailed = "SITE_CREATION_FAILED"
	ErrSiteUpdateFailed   = "SITE_UPDATE_FAILED"
	ErrSiteDeleteFailed   = "SITE_DELETE_FAILED"
	ErrSiteIDRequired     = "SITE_ID_REQUIRED"
	ErrAPIKeyGenFailed    = "API_KEY_GENERATION_FAILED"
	ErrBootstrapDone      = "BOOTSTRAP_ALREADY_DONE"
)

// Newsletter
const (
	ErrSubscribeFailed   = "SUBSCRIBE_FAILED"
	ErrUnsubscribeFailed = "UNSUBSCRIBE_FAILED"
	ErrIDGenFailed       = "ID_GENERATION_FAILED"
)

// Blog
const (
	ErrPostNotFound       = "POST_NOT_FOUND"
	ErrPostCreationFailed = "POST_CREATION_FAILED"
	ErrPostUpdateFailed   = "POST_UPDATE_FAILED"
	ErrPostDeleteFailed   = "POST_DELETE_FAILED"
	ErrPostPublishFailed  = "POST_PUBLISH_FAILED"
	ErrSlugExists         = "SLUG_EXISTS"
	ErrInvalidContent     = "INVALID_CONTENT"
	ErrInvalidDatetime    = "INVALID_DATETIME"
	ErrSchedulePastDate   = "SCHEDULE_PAST_DATE"
	ErrCommentFailed      = "COMMENT_FAILED"
	ErrDailyLimitReached  = "DAILY_LIMIT_REACHED"
	ErrAuthRequired       = "AUTHENTICATION_REQUIRED"
)

// Forum
const (
	ErrTopicNotFound       = "TOPIC_NOT_FOUND"
	ErrTopicCreationFailed = "TOPIC_CREATION_FAILED"
	ErrTopicClosed         = "TOPIC_CLOSED"
	ErrCategoryNotFound    = "CATEGORY_NOT_FOUND"
	ErrTagNotFound         = "TAG_NOT_FOUND"
	ErrTagExists           = "TAG_EXISTS"
)

// Document Access
const (
	ErrDocNotFound    = "DOCUMENT_NOT_FOUND"
	ErrAccessDenied   = "ACCESS_DENIED"
	ErrInvalidAPIKey  = "INVALID_API_KEY"
	ErrInvalidPath    = "INVALID_PATH"
	ErrPathRequired   = "PATH_REQUIRED"
	ErrRuleNotFound   = "RULE_NOT_FOUND"
	ErrRuleExists     = "RULE_EXISTS"
	ErrRuleIDRequired = "RULE_ID_REQUIRED"
)

// Upload
const (
	ErrUploadFailed     = "UPLOAD_FAILED"
	ErrUploadNotFound   = "UPLOAD_NOT_FOUND"
	ErrNoFileProvided   = "NO_FILE_PROVIDED"
	ErrStorageNotConfig = "STORAGE_NOT_CONFIGURED"
	ErrNotUploadOwner   = "NOT_UPLOAD_OWNER"
)

// Health
const (
	ErrDBUnreachable = "DATABASE_UNREACHABLE"
)

// respondError sends a JSON error response with a stable error code.
// Format: {"error": "human message", "code": "ERROR_CODE"}
func respondError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"error": message,
		"code":  code,
	})
}

// respondBadRequest is a shorthand for 400 errors.
func respondBadRequest(c *gin.Context, code string, message string) {
	respondError(c, http.StatusBadRequest, code, message)
}

// respondUnauthorized is a shorthand for 401 errors.
func respondUnauthorized(c *gin.Context, code string, message string) {
	respondError(c, http.StatusUnauthorized, code, message)
}

// respondNotFound is a shorthand for 404 errors.
func respondNotFound(c *gin.Context, code string, message string) {
	respondError(c, http.StatusNotFound, code, message)
}

// respondInternalError is a shorthand for 500 errors.
func respondInternalError(c *gin.Context, code string, message string) {
	respondError(c, http.StatusInternalServerError, code, message)
}
