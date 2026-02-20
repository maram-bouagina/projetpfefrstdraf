package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server
	ServerPort string
}

func Load() (Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %v", err)
	}

	// Extract PostgreSQL variables
	dbHost, err := extractEnv("DB_HOST")
	if err != nil {
		return Config{}, err
	}

	dbPort, err := extractEnv("DB_PORT")
	if err != nil {
		return Config{}, err
	}

	dbUser, err := extractEnv("DB_USER")
	if err != nil {
		return Config{}, err
	}

	dbPassword, err := extractEnv("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	dbName, err := extractEnv("DB_NAME")
	if err != nil {
		return Config{}, err
	}

	// SSL mode is optional (defaults to disable)
	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	// Server port
	serverPort, err := extractEnv("APP_PORT")
	if err != nil {
		return Config{}, err
	}

	return Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSSLMode:  dbSSLMode,
		ServerPort: serverPort,
	}, nil
}

func extractEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return val, nil
}

// GetDBConnectionString returns PostgreSQL connection string
func (c Config) GetDBConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// GetServerAddr returns server address with port
func (c Config) GetServerAddr() string {
	return ":" + c.ServerPort
}
