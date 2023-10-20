package config

import (
	"log"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Структура конфигурации
// Файл конфигурации должен распологаться в корне проекта и иметь имя config.yaml
type Config struct {
	Env     string `env:"APP_ENV" env-default:"local"`
	LogFile string `env:"APP_LOG_FILE" env-default:""`
	Server  Server
	Client  Client
	Loger   *slog.Logger
}

type Server struct {
	BindAddr    string        `env:"SERVER_BIND_ADDR" env-default:":8080"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"IDDLE" env-default:"30s"`
	SSLEnabled  bool          `env:"SSL_ENABLE" env-default:"false"`
	KeyChain    string        `env:"SSL_CHAIN" env-default:""`
	PrivateKey  string        `env:"SSL_KEY" env-default:""`
}

type Client struct {
	Timeout   time.Duration `env:"CLIENT_TIMEOUT" env-default:"5s"`
	UserAgent string        `env:"USER_AGENT" env-default:"Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko"`
}

func MustLoad() *Config {
	var instance Config
	if err := cleanenv.ReadEnv(&instance); err != nil {
		log.Fatalf("Ошибка чтения переменных окружения. %s", err.Error())
	}

	return &instance
}
