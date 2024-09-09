// Парсинг yml файла
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"` //env-default - значение по умолчанию, env-required: "true" - значение должно быть обязательно установлено
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string
	Password    string `env:"HTTP_SERVER_PASSWORD"`
}

// Must используется, когда вместо ошибки функция будет паниковать
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH") //Получаем путь из переменной окружения
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) { //Проверяем существование файла
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	cfg.HttpServer.User = os.Getenv("USER")
	if cfg.HttpServer.User == "" {
		log.Fatal("USER is not set")
	}
	cfg.HttpServer.Password = os.Getenv("HTTP_SERVER_PASSWORD")
	if cfg.HttpServer.Password == "" {
		log.Fatal("HTTP_SERVER_PASSWORD is not set")
	}

	return &cfg
}
