package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	addr   string
	logger *slog.Logger
}

func NewServer(addr string, logger *slog.Logger) *Server {
	return &Server{
		addr:   addr,
		logger: logger,
	}
}

func (s *Server) Start() (err error) {
	server := &http.Server{
		Addr:    s.addr,
		Handler: nil, // Используется http.DefaultServeMux
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server listen and serve error", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = s.Stop(server); err != nil {
		return
	}
	return
}

func (s *Server) Stop(server *http.Server) (err error) {
	s.logger.Info("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		s.logger.Error("Server shutdown failed", err)
		return
	}

	s.logger.Info(fmt.Sprintf("Server has done: %s", s.addr))
	return
}
