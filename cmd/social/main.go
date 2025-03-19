package main

import (
	"log/slog"
	"net/http"
	"os"

	"social-network/internal/config"
	"social-network/internal/handler"
	"social-network/internal/server"
	"social-network/internal/storage/postgres"
)

func main() {
	logger := setupLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", err)
		return
	}

	storage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		return
	}
	defer storage.Close()

	h := handler.New(logger, storage)
	http.HandleFunc("/login", h.MethodHandler(map[string]http.HandlerFunc{
		http.MethodPost: h.Login(),
	}))
	http.HandleFunc("/user/register", h.MethodHandler(map[string]http.HandlerFunc{
		http.MethodPost: h.Register(),
	}))
	http.HandleFunc("/user/get/", server.AuthMiddleware(
		h.MethodHandler(map[string]http.HandlerFunc{
			http.MethodGet: h.Users(),
		})),
	)

	srv := server.NewServer(cfg.HTTPAddr, logger)
	if err := srv.Start(); err != nil {
		return
	}
}

func setupLogger() (log *slog.Logger) {
	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return
}
