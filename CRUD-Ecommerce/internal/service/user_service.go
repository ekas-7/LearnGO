package service

import (
	"fmt"
	"time"

	"github.com/ekas-7/CRUD-Ecommerce/internal/model"
	"github.com/ekas-7/CRUD-Ecommerce/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req *model.UserRegisterRequest) (*model.UserResponse, error)
	Login(req *model.UserLoginRequest) (*model.LoginResponse, error)
	GetProfile(userID uuid.UUID) (*model.UserResponse, error)
	UpdateProfile(userID uuid.UUID, req *model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteAccount(userID uuid.UUID) error
	ValidateToken(tokenString string) (uuid.UUID, string, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewUserService(repo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *userService) Register(req *model.UserRegisterRequest) (*model.UserResponse, error) {
	// Check if user already exists
	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) Login(req *model.UserLoginRequest) (*model.LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	userResponse := user.ToResponse()
	return &model.LoginResponse{
		Token: token,
		User:  userResponse,
	}, nil
}

func (s *userService) GetProfile(userID uuid.UUID) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) DeleteAccount(userID uuid.UUID) error {
	return s.repo.Delete(userID)
}

func (s *userService) generateToken(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *userService) ValidateToken(tokenString string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, "", fmt.Errorf("invalid user_id in token")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, "", fmt.Errorf("invalid user_id format: %w", err)
		}

		role, ok := claims["role"].(string)
		if !ok {
			role = "user"
		}

		return userID, role, nil
	}

	return uuid.Nil, "", fmt.Errorf("invalid token claims")
}
