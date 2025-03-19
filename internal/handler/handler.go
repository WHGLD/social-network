package handler

import (
	"log/slog"
	"net/http"

	model "social-network/internal/models"
)

type Storage interface {
	CreateUser(user *model.User) error
	GetUserByID(userID string) (*model.User, error)
}

type Handler struct {
	logger  *slog.Logger
	storage Storage
}

// MethodHandler возвращает функцию-обработчик, которая выбирает нужную функцию по методу запроса.
func (h *Handler) MethodHandler(methodHandlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := methodHandlers[r.Method]; ok {
			handler(w, r)
		} else {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func New(logger *slog.Logger, storage Storage) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
	}
}
