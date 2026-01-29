// Пакет config отвечает за загрузку конфигурации из переменных окружения.
package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит все параметры приложения.
type Config struct {
	// PostgreSQL
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresPort     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUser     string `envconfig:"POSTGRES_USER" default:"app_user"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"secure_password"`
	PostgresDB       string `envconfig:"POSTGRES_DB" default:"app_db"`

	// MongoDB
	MongoURI string `envconfig:"MONGO_URI" default:"mongodb://localhost:27017"`
	MongoDB  string `envconfig:"MONGO_DB" default:"app_db"`

	// JWT
	JWTSecret     string        `envconfig:"JWT_SECRET" default:"my_super_secret_key_123!"`
	JWTExpiresIn  time.Duration `envconfig:"JWT_EXPIRES_IN" default:"24h"`

	// Сервер
	ServerPort string `envconfig:"SERVER_PORT" default:"8080"`
}

// Load читает переменные окружения и заполняет структуру Config.
// Возвращает ошибку, если переменные имеют неверный формат.
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}


