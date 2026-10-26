package handler

import (
	"Museum/internal/auth"
	"Museum/internal/models"
	"Museum/internal/service"
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=3"`
	FullName string          `json:"full_name" binding:"required,min=6"`
	Role     models.UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Role == "" ||
		!slices.Contains([]models.UserRole{models.RoleStaff, models.RoleAdmin}, req.Role) {
		req.Role = models.RoleStaff
	}
	u, token, err := h.svc.Register(c.Request.Context(), req.Email, req.Password, req.FullName, req.Role)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, AuthResponse{User: u, Token: token})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, token, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, AuthResponse{User: u, Token: token})
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=6"`
}

func (h *UserHandler) UpdateInfo(c *gin.Context) {
	uid := c.GetInt64(auth.UserIdKey)
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.svc.Update(c.Request.Context(), uid, req.FullName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	uid := c.GetInt64(auth.UserIdKey)
	u, err := h.svc.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}
