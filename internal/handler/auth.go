package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Dided08/Calculator/internal/models"
	"github.com/Dided08/Calculator/internal/token"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Logger     *zap.Logger
	Storage    UserStorage // интерфейс репозитория
	JWTSecret  string
	TokenTTL   time.Duration
}

type UserStorage interface {
	CreateUser(username, password string) error
	AuthenticateUser(username, password string) (int, error)
}

// RegisterHandler обрабатывает регистрацию пользователя
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Имя пользователя и пароль обязательны", http.StatusUnprocessableEntity)
		return
	}

	err := h.Storage.CreateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Пользователь уже существует или ошибка сервера", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// LoginHandler обрабатывает вход пользователя
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	userID, err := h.Storage.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	tokenStr, err := token.GenerateToken(userID, h.JWTSecret, h.TokenTTL)
	if err != nil {
		h.Logger.Error("ошибка генерации токена", zap.Error(err))
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}

	resp := models.TokenResponse{Token: tokenStr}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}