package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) Users() http.HandlerFunc {
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
