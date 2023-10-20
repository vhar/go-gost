package apiserver

import (
	"go-stribog/internal/apiserver/handlers"
	"go-stribog/internal/config"
	"log/slog"

	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/cors"
)

type APIServer struct {
	config *config.Config
	logger *slog.Logger
	router *chi.Mux
}

func New(config *config.Config) *APIServer {
	return &APIServer{
		config: config,
		logger: config.Loger,
		router: chi.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.configureRoutes()

	serv := &http.Server{
		Addr:         s.config.Server.BindAddr,
		Handler:      s.router,
		ReadTimeout:  s.config.Server.Timeout,
		WriteTimeout: s.config.Server.Timeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	if s.config.Server.SSLEnabled == true {
		s.logger.Info("Запускаем API Server с SSL...")

		if err := serv.ListenAndServeTLS(s.config.Server.KeyChain, s.config.Server.PrivateKey); err != nil {
			return fmt.Errorf("Ощибка запуска api сервера. %s", err.Error())
		}
	} else {
		s.logger.Info("Запускаем API Server без SSL...")

		if err := serv.ListenAndServe(); err != nil {
			return fmt.Errorf("Ощибка запуска api сервера. %s", err.Error())
		}
	}

	return nil
}

func (s *APIServer) configureRoutes() {
	s.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	h := handlers.New(s.config.Client, s.logger)

	s.router.Post("/vfile", h.VerifyFile)
	s.router.Post("/vurl", h.VerifyURL)
}
