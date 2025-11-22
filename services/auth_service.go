package services

import (
	"errors"
	"googleforms/internal/dto"
	"googleforms/internal/models"
	"googleforms/repositories"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// if !user.IsActive {
	// 	return nil, errors.New("account is deactivated")
	// }

	// generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			TenantID:  user.TenantID,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {

	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		TenantID:     req.TenantID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         "editor",
		// IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		TenantID:  user.TenantID,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID := uint(claims["user_id"].(float64))
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"email":     user.Email,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 24 ساعة
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
