package cmd

import (
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

type ServerInterface interface {
	Start()
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

func (s *Server) Start() {
	s.log.Info("Starting the server")
	s.app = fiber.New(fiber.Config{
		IdleTimeout: 60 * 60,

		ReadTimeout:       60 * 60,
		ReduceMemoryUsage: true,
		WriteTimeout:      60 * 60,
		WriteBufferSize:   1024 * 1024,
	})

	s.setupMiddleware()
	s.setupRoutes()
	if err := s.app.Listen(":8545"); err != nil {
		s.log.Panic(err.Error())
	}

	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt)
	<-osChannel
	if err := s.app.Shutdown(); err != nil {
		s.log.Panic(err.Error())
	}
}

func (s *Server) setupMiddleware() {
	s.log.Info("Setting up the middleware")
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))
}

func (s *Server) setupRoutes() {
	s.log.Info("Setting up the routes")
	s.app.Post("/submit", s.handler.SubmitTx)
	s.app.Get("/tx/:txHash", s.handler.GetTx)
}
