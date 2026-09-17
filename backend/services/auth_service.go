package services

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "atk-backend/models"
    "atk-backend/repositories"
    "atk-backend/utils"
)

type AuthService struct {
    userRepo *repositories.UserRepository
    roleRepo *repositories.RoleRepository
}

func NewAuthService() *AuthService {
    return &AuthService{
        userRepo: repositories.NewUserRepository(),
        roleRepo: repositories.NewRoleRepository(),
    }
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username    string `json:"username" binding:"required"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=6"`
    FullName    string `json:"full_name" binding:"required"`
    PhoneNumber string `json:"phone_number"`
}

type LoginResponse struct {
    Token string      `json:"token"`
    User  models.User `json:"user"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
    user, err := s.userRepo.FindByUsername(req.Username)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    if user.Status != "active" {
        return nil, errors.New("account is not active")
    }
    
    if !utils.CheckPasswordHash(req.Password, user.Password) {
        return nil, errors.New("invalid credentials")
    }
    
    // Update last login
    now := time.Now()
    user.LastLogin = &now
    s.userRepo.Update(user)
    
    // Get primary role
    role := ""
    if len(user.Roles) > 0 {
        role = user.Roles[0].Name
    }
    
    token, err := utils.GenerateToken(user.ID, user.Username, user.Email, role)
    if err != nil {
        return nil, err
    }
    
    return &LoginResponse{
        Token: token,
        User:  *user,
    }, nil
}

func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
    // Check if username exists
    _, err := s.userRepo.FindByUsername(req.Username)
    if err == nil {
        return nil, errors.New("username already exists")
    }
    
    // Check if email exists
    _, err = s.userRepo.FindByEmail(req.Email)
    if err == nil {
        return nil, errors.New("email already exists")
    }
    
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        UUID:        uuid.New().String(),
        Username:    req.Username,
        Email:       req.Email,
        Password:    hashedPassword,
        FullName:    req.FullName,
        PhoneNumber: req.PhoneNumber,
        Status:      "active",
    }
    
    err = s.userRepo.Create(user)
    if err != nil {
        return nil, err
    }
    
    // Assign default role (staff)
    role, err := s.roleRepo.FindByName("staff")
    if err == nil {
        s.userRepo.AssignRole(user.ID, role.ID)
    }
    
    return user, nil
}

func (s *AuthService) ValidateToken(token string) (*utils.Claims, error) {
    return utils.ValidateToken(token)
}