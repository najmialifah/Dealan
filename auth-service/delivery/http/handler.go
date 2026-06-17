package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/najmialifah/Dealan/auth-service/domain"
	"github.com/najmialifah/Dealan/auth-service/service"
)

// AuthHandler adalah controller untuk REST API Autentikasi menggunakan Gin
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler membuat instance baru dari AuthHandler
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register menangani registrasi akun baru
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login menangani autentikasi akun
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Validate memvalidasi token JWT dari header Authorization
func (h *AuthHandler) Validate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header tidak ditemukan"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid (gunakan Bearer <token>)"})
		return
	}

	cred, err := h.svc.ValidateToken(c.Request.Context(), tokenParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cred)
}
