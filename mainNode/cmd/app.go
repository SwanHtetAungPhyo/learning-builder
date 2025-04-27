package cmd

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerInterface interface {
	Start() error
	setupMiddleware()
	setupRoutes()
}

type Server struct {
	log     *logrus.Logger
	app     *fiber.App
	handler handler.Handler
}

var _ ServerInterface = (*Server)(nil)

func NewServer(log *logrus.Logger, handle handler.Handler) *Server {
	return &Server{
		log:     log,
		handler: handle,
	}
}

func (s *Server) Start() error {
	s.log.Info("Initializing server...")

	// Create Fiber app with sensible defaults
	s.app = fiber.New(fiber.Config{
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			s.log.Errorf("Request error: %v", err)
			return fiber.DefaultErrorHandler(ctx, err)
		},
	})

	// Setup middleware and routes
	s.setupMiddleware()
	s.setupRoutes()

	// Create channel for shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start server in separate goroutine
	serverErr := make(chan error, 1)
	go func() {
		s.log.Info("Starting server on :8545")
		if err := s.app.Listen(":8545"); err != nil {
			serverErr <- fmt.Errorf("server failed: %w", err)
		}
	}()

	// Wait for either shutdown signal or server error
	select {
	case sig := <-shutdownChan:
		s.log.Infof("Received %v signal, shutting down...", sig)
	case err := <-serverErr:
		s.log.Errorf("Server error: %v", err)
		return err
	}

	s.log.Info("Initiating graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	s.log.Info("Server shutdown complete")
	return nil
}

func (s *Server) setupMiddleware() {
	s.log.Info("Setting up middleware...")

	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// CORS configuration
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept",
		ExposeHeaders: "Content-Length",
		MaxAge:        86400,
	}))
}

func (s *Server) setupRoutes() {
	s.log.Info("Setting up routes...")

	s.app.Post("/submit", s.handler.SubmitTx)
	s.app.Get("/tx/:txHash", s.handler.GetTx)
	s.app.Post("/blocks-propose", s.handler.BlockAddition)

	// Health check endpoint
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})
}
