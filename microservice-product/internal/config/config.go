package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	/*holding the configuration values*/
	// PostgreSQL
	DBHost     string // win 9aad yrunni fi docker db-app wella localhost
	DBPort     string //port mtaa postgres
	DBUser     string //postgres user mtaa bdd
	DBPassword string // password mtaa  bdd
	DBName     string //nom mtaa bdd
	DBSSLMode  string //  encrypt mtaa connection disabled wella disabled par défaut

	// Server: 8000 par exemple
	ServerPort string
}

func Load() (Config, error) {
	// Load .env file //// ta9rah
	if err := godotenv.Load(); err != nil {
		//trajaa struct faraght wl erreur
		return Config{}, fmt.Errorf("error loading .env file: %v", err)
	}

	// tkhrj PostgreSQL variables
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

	// hia chtaaml: lconnection yrudha chifrée idhe keni vide urudha diabled mahich dechifrée
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

// lfunc illi staamlneha
func extractEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return val, nil
}

// nistha9ou ha string bech nconnectiw mawjouda fi connection.go
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
