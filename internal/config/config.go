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
	Env     string  `yaml:"env"`
	LogFile string  `yaml:"log_file"`
	Server  *Server `yaml:"server"`
	Client  *Client `yaml:"client"`
	Loger   *slog.Logger
}

type Server struct {
	BindAddr    string        `yaml:"bind_addr"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	SSLEnabled  bool          `yaml:"ssl_enabled"`
	KeyChain    string        `yaml:"ssl_chain"`
	PrivateKey  string        `yaml:"ssl_key"`
}

type Client struct {
	Timeout   time.Duration `yaml:"timeout"`
	UserAgent string        `yaml:"user_agent"`
}

func MustLoad() *Config {
	var instance Config
	if err := cleanenv.ReadConfig("config.yaml", &instance); err != nil {
		log.Fatalf("Ошибка чтения файла конфигурации. %s", err.Error())
	}

	return &instance
}
