package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/studiowebux/minimaldoc/internal/server/auth"
)

// User management handlers

func (r *Router) listUsers(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := r.db.ListUsers(c.Request.Context(), siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	// Convert to response format (exclude password hash)
	type userResponse struct {
		ID            string  `json:"id"`
		Email         string  `json:"email"`
		Role          string  `json:"role"`
		Name          *string `json:"name,omitempty"`
		EmailVerified bool    `json:"email_verified"`
		CreatedAt     string  `json:"created_at"`
		LastLoginAt   *string `json:"last_login_at,omitempty"`
	}

	result := make([]userResponse, len(users))
	for i, u := range users {
		result[i] = userResponse{
			ID:            u.ID,
			Email:         u.Email,
			Role:          u.Role,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt,
		}
		if u.Name.Valid {
			result[i].Name = &u.Name.String
		}
		if u.LastLoginAt.Valid {
			result[i].LastLoginAt = &u.LastLoginAt.String
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": result})
}

func (r *Router) createUser(c *gin.Context) {
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"required,oneof=admin editor author viewer"`
		Name     string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	existing, err := r.db.GetUserByEmail(c.Request.Context(), siteID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password, r.config.Auth.BCryptCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Generate user ID
	userID := uuid.New().String()

	// Create user
	user, err := r.db.CreateUser(c.Request.Context(), userID, siteID, req.Email, passwordHash, req.Role, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (r *Router) getUser(c *gin.Context) {
	id := c.Param("id")

	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := r.db.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check site scope
	if user.SiteID != siteID {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	response := gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"role":           user.Role,
		"email_verified": user.EmailVerified,
		"created_at":     user.CreatedAt,
		"updated_at":     user.UpdatedAt,
	}
	if user.Name.Valid {
		response["name"] = user.Name.String
	}
	if user.LastLoginAt.Valid {
		response["last_login_at"] = user.LastLoginAt.String
	}
	if user.OAuthProvider.Valid {
		response["oauth_provider"] = user.OAuthProvider.String
	}

	c.JSON(http.StatusOK, response)
}

func (r *Router) updateUser(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get existing user
	user, err := r.db.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if user == nil || user.SiteID != siteID {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Role     string `json:"role" binding:"omitempty,oneof=admin editor author viewer"`
		Name     string `json:"name"`
		Password string `json:"password" binding:"omitempty,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent self-privilege escalation: users cannot change their own role
	claims, err := getUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if claims.UserID == id && req.Role != "" && req.Role != user.Role {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot change your own role"})
		return
	}

	// Use existing values if not provided
	email := user.Email
	if req.Email != "" {
		email = req.Email
	}
	role := user.Role
	if req.Role != "" {
		role = req.Role
	}
	name := user.Name.String
	if req.Name != "" {
		name = req.Name
	}

	// Update user
	if err := r.db.UpdateUser(c.Request.Context(), id, email, role, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Update password if provided
	if req.Password != "" {
		passwordHash, err := auth.HashPassword(req.Password, r.config.Auth.BCryptCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		if err := r.db.UpdateUserPassword(c.Request.Context(), id, passwordHash); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (r *Router) deleteUser(c *gin.Context) {
	id := c.Param("id")
	siteID, err := getSiteID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get user to verify ownership
	user, err := r.db.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if user == nil || user.SiteID != siteID {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Prevent deleting self
	claims, err := getUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if claims.UserID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	// Delete user sessions first
	if err := r.db.DeleteUserSessions(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user sessions"})
		return
	}

	// Delete user
	if err := r.db.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
