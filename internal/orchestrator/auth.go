package orchestrator

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dided08/Calculator/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService предоставляет методы для регистрации и аутентификации пользователей.
type AuthService struct {
	DB *sql.DB
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{DB: db}
}

// Register регистрирует нового пользователя.
func (a *AuthService) Register(username, password string) error {
	// Проверка на дубликат
	var exists bool
	err := a.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки пользователя: %w", err)
	}
	if exists {
		return errors.New("пользователь уже существует")
	}

	// Хеширование пароля
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	// Вставка в БД
	_, err = a.DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hashed)
	if err != nil {
		return fmt.Errorf("ошибка вставки пользователя: %w", err)
	}

	return nil
}

// Login аутентифицирует пользователя и возвращает его ID.
func (a *AuthService) Login(username, password string) (int, error) {
	var user models.User

	err := a.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return 0, errors.New("неверный логин или пароль")
	} else if err != nil {
		return 0, fmt.Errorf("ошибка запроса пользователя: %w", err)
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return 0, errors.New("неверный логин или пароль")
	}

	return user.ID, nil
}