package handlers

import (
    "github.com/gin-gonic/gin"
    "atk-backend/services"
    "atk-backend/utils"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{
        authService: services.NewAuthService(),
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req services.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequestResponse(c, "Invalid request", err.Error())
        return
    }
    
    resp, err := h.authService.Login(req)
    if err != nil {
        utils.UnauthorizedResponse(c, err.Error())
        return
    }
    
    utils.SuccessResponse(c, "Login successful", resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req services.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequestResponse(c, "Invalid request", err.Error())
        return
    }
    
    user, err := h.authService.Register(req)
    if err != nil {
        utils.BadRequestResponse(c, err.Error(), nil)
        return
    }
    
    utils.CreatedResponse(c, "User registered successfully", user)
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
    userID := c.GetUint("user_id")
    username := c.GetString("username")
    
    utils.SuccessResponse(c, "Token is valid", gin.H{
        "user_id":  userID,
        "username": username,
    })
}