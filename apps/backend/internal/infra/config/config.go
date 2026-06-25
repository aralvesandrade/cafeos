package config

import "os"

type Config struct {
	ServerPort  string
	DatabaseURL string
	RedisURL    string
	RabbitMQURL string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "5001"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://cafeos:cafeos@localhost:5432/cafeos?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://cafeos:cafeos@localhost:5672/"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
