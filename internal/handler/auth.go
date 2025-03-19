package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"social-network/internal/common"
	model "social-network/internal/models"
)

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			creds struct {
				UserID   string `json:"user_id"`
				Password string `json:"password"`
			}
			isPwValid bool
			token     string
			user      *model.User
			err       error
		)
		if err = json.NewDecoder(r.Body).Decode(&creds); err != nil {
			h.logger.Error("Invalid request body", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err = h.storage.GetUserByID(creds.UserID)
		if err != nil {
			h.logger.Error("User not found", err)
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if isPwValid, err = matchPassword(user.PasswordHash, creds.Password); isPwValid == false || err != nil {
			if err != nil {
				h.logger.Error("User password is not correct. ERR", err)
			}
			http.Error(w, "User password is not correct", http.StatusUnauthorized)
			return
		}

		token, err = getToken(user.ID)
		if err != nil {
			h.logger.Error("Failed to generate token", err)
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user model.User
			err  error
		)

		if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
			h.logger.Error("Invalid request data", err)
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		if user.Birthday != nil {
			if _, err = time.Parse("2006-01-02", *user.Birthday); err != nil {
				http.Error(w, "Invalid request data (birthday)", http.StatusBadRequest)
				return
			}
		}

		if user.Sex != nil {
			if !inArrayStr(*user.Sex, []string{model.GenderFemale, model.GenderMale}) {
				http.Error(w, "Invalid request data (sex)", http.StatusBadRequest)
				return
			}
		}

		user.ID = uuid.NewString()
		user.PasswordHash, err = hashPassword(user.Password)
		if err != nil {
			h.logger.Error("Failed to hash password", err)
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		if err = h.storage.CreateUser(&user); err != nil {
			h.logger.Error("Failed to create user", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"user_id": user.ID})
	}
}

func inArrayStr(value string, array []string) (result bool) {
	for i := 0; i < len(array); i++ {
		if value == array[i] {
			result = true
			return
		}
	}
	return
}

func matchPassword(hashedPw, pw string) (result bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw)); err != nil {
		result = false
		return
	}

	result = true
	return
}

func getToken(userID string) (authToken string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &common.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	authToken, err = token.SignedString(common.JwtKey)
	if err != nil {
		return "", err
	}

	return
}

func hashPassword(pw string) (hash string, err error) {
	generatedHash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash = string(generatedHash)
	return
}
