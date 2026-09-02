package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	AppPort    string
	DBDriver   string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	JWTSecret  string
	JWTExpires int
	UploadPath string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	jwtExpires, err := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "72"))
	if err != nil {
		jwtExpires = 72
	}

	return &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBDriver:   getEnv("DB_DRIVER", "mysql"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "iot_rd_master"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		JWTSecret:  getEnv("JWT_SECRET", "iot-rd-control-center-secret-jwt-key-2026"),
		JWTExpires: jwtExpires,
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
