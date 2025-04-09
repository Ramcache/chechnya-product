package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Регистрация нового пользователя
func (s *UserService) Register(username, password string) error {
	// Проверим, существует ли уже пользователь
	existing, _ := s.repo.GetByUsername(username)
	if existing != nil {
		return errors.New("пользователь уже существует")
	}

	// Хешируем пароль
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Сохраняем пользователя
	user := &models.User{
		Username: username,
		Password: string(hashed),
		Role:     "user", // по умолчанию
	}
	return s.repo.Create(user)
}

// Аутентификация (логин)
func (s *UserService) Login(username, password string) (*models.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("неверный пароль")
	}

	return user, nil
}

// Получить пользователя по ID
func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}
