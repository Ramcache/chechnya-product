package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register(req RegisterRequest) (*models.User, error)
	LoginWithUser(req LoginRequest) (*models.User, string, error)
	GetByID(userID int) (*models.User, error)
	GetByOwnerID(ownerID string) (*models.User, error)
	TransferCart(oldOwnerID, newOwnerID string) error
	CreateByPhone(phone string) (*models.User, string, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type UserService struct {
	repo        repositories.UserRepository
	jwt         utils.JWTManagerInterface
	cartService CartServiceInterface
}

func NewUserService(repo repositories.UserRepository, jwt utils.JWTManagerInterface, cart CartServiceInterface) *UserService {
	return &UserService{repo: repo, jwt: jwt, cartService: cart}
}

// Данные для регистрации пользователя
type RegisterRequest struct {
	Phone    string
	Password string
	Username string
	Email    *string
	OwnerID  string
}

// Данные для входа пользователя
type LoginRequest struct {
	Identifier string // телефон, почта или имя
	Password   string
}

// Регистрация пользователя по телефону
func (s *UserService) Register(req RegisterRequest) (*models.User, error) {
	// Валидация (например, длина пароля, пустой username и т.д.)
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Проверки на существование (ускоряет ответ)
	existingPhone, _ := s.repo.GetByPhone(req.Phone)
	if existingPhone != nil {
		return nil, errors.New("Пользователь с таким номером уже существует")
	}

	if req.Email != nil && *req.Email != "" {
		existingEmail, _ := s.repo.GetByEmail(*req.Email)
		if existingEmail != nil {
			return nil, errors.New("Пользователь с таким email уже существует")
		}
	}

	existingUsername, _ := s.repo.GetByUsername(req.Username)
	if existingUsername != nil {
		return nil, errors.New("Пользователь с таким именем уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать хэш пароля: %w", err)
	}

	user := &models.User{
		Username:     req.Username,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         models.UserRoleUser,
		IsVerified:   true,
		OwnerID:      "user_" + uuid.New().String(), // всегда новый уникальный owner_id!
		CreatedAt:    time.Now(),
	}

	// Пытаемся создать пользователя — ловим уникальные ограничения
	if err := s.repo.Create(user); err != nil {
		fmt.Printf("Ошибка при создании пользователя: %+v\n", err)
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			switch {
			case strings.Contains(pgErr.Message, "users_email_key"):
				return nil, errors.New("Пользователь с таким email уже зарегистрирован")
			case strings.Contains(pgErr.Message, "users_phone_key"):
				return nil, errors.New("Пользователь с таким номером уже зарегистрирован")
			case strings.Contains(pgErr.Message, "users_username_key"):
				return nil, errors.New("Пользователь с таким именем уже зарегистрирован")
			}
		}
		return nil, errors.New("Не удалось создать пользователя")
	}

	return user, nil
}

// Аутентификация пользователя (вход)
// Возвращает пользователя и JWT токен
func (s *UserService) LoginWithUser(req LoginRequest) (*models.User, string, error) {
	if err := utils.ValidateIdentifier(req.Identifier); err != nil {
		return nil, "", err
	}

	// Определяем тип идентификатора (телефон, e-mail или username)
	var user *models.User
	var err error
	switch {
	case strings.Contains(req.Identifier, "@"):
		user, err = s.repo.GetByEmail(req.Identifier)
	case strings.HasPrefix(req.Identifier, "+"):
		user, err = s.repo.GetByPhone(req.Identifier)
	default:
		user, err = s.repo.GetByUsername(req.Identifier)
	}
	if err != nil || user == nil {
		return nil, "", errors.New("Неверные данные для входа")
	}

	if !user.IsVerified {
		return nil, "", errors.New("Сначала подтвердите номер телефона")
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, "", errors.New("Неверный пароль")
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("Ошибка генерации токена: %w", err)
	}

	return user, token, nil
}

// Получить пользователя по ID
func (s *UserService) GetByID(userID int) (*models.User, error) {
	return s.repo.GetByID(userID)
}

// Получить пользователя по owner_id
func (s *UserService) GetByOwnerID(ownerID string) (*models.User, error) {
	return s.repo.GetByOwnerID(ownerID)
}

// Перенос корзины при смене пользователя
func (s *UserService) TransferCart(oldOwnerID, newOwnerID string) error {
	return s.cartService.TransferCart(oldOwnerID, newOwnerID)
}

// Валидация данных регистрации
func validateRegisterRequest(req RegisterRequest) error {
	if err := utils.ValidatePhone(req.Phone); err != nil {
		return err
	}
	if len(req.Password) < 6 {
		return errors.New("Пароль должен быть не менее 6 символов")
	}
	return nil
}

// Проверка пароля
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Создание пользователя только по номеру телефона (для админов)
func (s *UserService) CreateByPhone(phone string) (*models.User, string, error) {
	if err := utils.ValidatePhone(phone); err != nil {
		return nil, "", err
	}

	existing, _ := s.repo.GetByPhone(phone)
	if existing != nil {
		return nil, "", errors.New("Пользователь уже существует")
	}

	password := utils.GeneratePassword(8)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Phone:        phone,
		PasswordHash: string(hash),
		Username:     "user_" + phone[len(phone)-4:],
		Role:         models.UserRoleUser,
		IsVerified:   true,
		OwnerID:      "user_" + utils.GenerateShortID(),
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, "", err
	}

	return user, password, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}
