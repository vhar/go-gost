package handlers

import (
	"go-gost/internal/config"
	"go-gost/internal/lib/signature"
	"log/slog"
)

type Handler struct {
	config config.Client
	logger *slog.Logger
}

// Структура отчета о проверке ЭЦП
type report struct {
	Signer      signature.Signer
	Certificate signature.Certificate
	SigningTime string
	Validity    bool
}

func New(config config.Client, logger *slog.Logger) *Handler {
	return &Handler{
		config: config,
		logger: logger,
	}
}
