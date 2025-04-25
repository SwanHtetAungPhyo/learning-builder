package main

import (
	"fmt"

	"github.com/SwanHtetAungPhyo/learning/mainNode/cmd"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/handler"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

type Config struct {
	Port        string
	RabbitMQURL string
}

var config *Config

func NewConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Warnf(err.Error())
	}
	port := os.Getenv("PORT")
	if port == " " {
		port = "8545"
	}
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	config = &Config{
		Port:        port,
		RabbitMQURL: rabbitMQURL,
	}
}
func init() {
	NewConfig()
}
func main() {
	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filePlusLine := fmt.Sprintf("%s:%d", f.File, f.Line)
			return f.Function, filePlusLine
		},
		PrettyPrint: true,
	}

	server := cmd.NewServer(logger, handler.NewImpl(logger))

	server.Start()

}
