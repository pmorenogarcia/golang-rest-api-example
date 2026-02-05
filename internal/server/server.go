package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/polgarcia/golang-rest-api/internal/config"
	"github.com/polgarcia/golang-rest-api/internal/handler"
	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

// New creates a new HTTP server
func New(cfg *config.Config, h *handler.Handler, log *logger.Logger) *Server {
	// Setup routes
	router := SetupRoutes(h, log, cfg.CORS.AllowedOrigins)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		logger:     log,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		zap.String("addr", s.httpServer.Addr),
	)

	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	s.logger.Info("Server started successfully",
		zap.String("addr", s.httpServer.Addr),
	)

	// Wait for interrupt signal
	return s.waitForShutdown()
}

// waitForShutdown waits for an interrupt signal and performs graceful shutdown
func (s *Server) waitForShutdown() error {
	// Create channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	sig := <-quit
	s.logger.Info("Received shutdown signal",
		zap.String("signal", sig.String()),
	)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	s.logger.Info("Shutting down server gracefully...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Server shutdown failed", zap.Error(err))
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info("Server stopped gracefully")
	return nil
}

// Stop stops the HTTP server immediately
func (s *Server) Stop() error {
	return s.httpServer.Close()
}
