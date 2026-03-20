package dto

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// TokenResponse represents the response containing a JWT token
type TokenResponse struct {
	Token string `json:"token"`
}

// RefreshRequest represents the request body for refreshing a token
type RefreshRequest struct {
	Username string `json:"username"`
}