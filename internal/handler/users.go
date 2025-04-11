package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	model "social-network/internal/models"
)

func (h *Handler) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/user/get/")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		user, err := h.storage.GetUserByID(id)
		if err != nil {
			h.logger.Error("User not found", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user.ToResponse())
	}
}

func (h *Handler) UsersSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		firstName := query.Get("first_name")
		secondName := query.Get("last_name")
		if (firstName == "" || len(firstName) > 20) && (secondName == "" || len(secondName) > 20) {
			http.Error(w, "Query params are not valid", http.StatusBadRequest)
			return
		}

		users, err := h.storage.UsersSearch(firstName, secondName)
		if err != nil {
			h.logger.Error("Users not found", err)
			http.Error(w, "Users not found", http.StatusNotFound)
			return
		}
		if len(users) == 0 {
			h.logger.Error("Users not found. len(users): 0")
			http.Error(w, "Users not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.UsersTransform(users))
	}
}
