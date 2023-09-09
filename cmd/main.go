package main

import (
	"fmt"
	"go-stribog/internal/apiserver"
	"go-stribog/internal/config"
	"log"
	"log/slog"
	"os"
)

var (
	logOutput *os.File
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	if err := configureLogger(config); err != nil {
		log.Fatal(err.Error())
	}

	s := apiserver.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

// Конфигурация логирования
func configureLogger(config *config.Config) (err error) {
	logOutput = os.Stdout
	if config.LogFile != "" {
		logOutput, err = os.OpenFile(fmt.Sprintf("%s", config.LogFile), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("Невозможно создать файл %s", config.LogFile)
		}
	}

	switch config.Env {
	case envLocal:
		config.Loger = slog.New(
			slog.NewTextHandler(logOutput, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		config.Loger = slog.New(
			slog.NewJSONHandler(logOutput, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		config.Loger = slog.New(
			slog.NewJSONHandler(logOutput, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return nil
}
