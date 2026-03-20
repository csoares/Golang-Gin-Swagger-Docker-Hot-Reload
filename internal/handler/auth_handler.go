package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"projetoapi/internal/dto"
	"projetoapi/internal/infrastructure/jwt"
	"projetoapi/internal/service"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	userService service.UserService
	jwtService  jwt.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService service.UserService, jwtService jwt.Service) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Login godoc
// @Summary Realizar autenticação
// @Description Autentica o utilizador e gera o token para os próximos acessos
// @Accept json
// @Produce json
// @Router /auth/login [post]
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad request!"})
		return
	}

	user, err := h.userService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid credentials!"})
		return
	}

	token, err := h.jwtService.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Success!", "token": token})
}

// Register godoc
// @Summary Realizar registro
// @Description Regista um utilizador
// @Accept json
// @Produce json
// @Router /auth/register [post]
// @Param user body dto.RegisterRequest true "Registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad request!"})
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Username already exists!"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Success!", "userId": user.ID})
}

// RefreshToken godoc
// @Summary Atualiza token de autenticação
// @Description Atualiza o token de autenticação do usuário
// @Accept json
// @Produce json
// @Router /auth/refresh_token [put]
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	username, err := h.jwtService.GetUsernameFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Not authorized"})
		return
	}

	token, err := h.jwtService.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Token updated successfully!", "token": token})
}
