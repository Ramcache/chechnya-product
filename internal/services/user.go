package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register(req RegisterRequest) (*models.User, error)
	Login(req LoginRequest) (string, error)
	GetByID(userID int) (*models.User, error)
	GetByOwnerID(ownerID string) (*models.User, error)
	TransferCart(oldOwnerID, newOwnerID string) error
}

type UserService struct {
	repo        repositories.UserRepository
	jwt         utils.JWTManagerInterface
	cartService CartServiceInterface
}

func NewUserService(repo repositories.UserRepository, jwt utils.JWTManagerInterface, cart CartServiceInterface) *UserService {
	return &UserService{repo: repo, jwt: jwt, cartService: cart}
}

// Запрос на регистрацию
type RegisterRequest struct {
	Phone    string
	Password string
	Username string
	Email    *string
	OwnerID  string
}

// Запрос на вход
type LoginRequest struct {
	Identifier string // username, phone или email
	Password   string
}

// Регистрация пользователя по номеру телефона
func (s *UserService) Register(req RegisterRequest) (*models.User, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	existing, _ := s.repo.GetByPhone(req.Phone)
	if existing != nil {
		return nil, errors.New("user with this phone already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     req.Username,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         models.UserRoleUser,
		IsVerified:   false,
		OwnerID:      req.OwnerID,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Авторизация по номеру телефона
func (s *UserService) Login(req LoginRequest) (string, error) {
	var user *models.User
	var err error

	// Определяем тип идентификатора
	switch {
	case isValidEmail(req.Identifier):
		user, err = s.repo.GetByEmail(req.Identifier)
	case isValidPhone(req.Identifier):
		user, err = s.repo.GetByPhone(req.Identifier)
	default:
		user, err = s.repo.GetByUsername(req.Identifier)
	}

	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return "", errors.New("phone not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// Получить пользователя по owner_id
func (s *UserService) GetByOwnerID(ownerID string) (*models.User, error) {
	return s.repo.GetByOwnerID(ownerID)
}

// Получить пользователя по ID
func (s *UserService) GetByID(userID int) (*models.User, error) {
	return s.repo.GetByID(userID)
}

func (s *UserService) TransferCart(oldOwnerID, newOwnerID string) error {
	return s.cartService.TransferCart(oldOwnerID, newOwnerID)
}

// Валидация регистрации
func validateRegisterRequest(req RegisterRequest) error {
	if strings.TrimSpace(req.Phone) == "" || len(req.Phone) < 10 {
		return errors.New("invalid phone number")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func isValidEmail(s string) bool {
	return strings.Contains(s, "@")
}

func isValidPhone(s string) bool {
	// Простейшая проверка: начинается с "+" и длина от 10
	return strings.HasPrefix(s, "+") && len(s) >= 10
}
