package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	DBUser         string
	DBPassword     string
	DBName         string
	DBHost         string
	DBPort         string
	IPInfoToken    string
	JWTSecret      string
	JWTExpiration  int
	RedisURL       string
	MailerHost     string
	MailerPort     int
	MailerUsername string
	MailerPassword string
	EncryptionKey  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "12345678"),
		DBName:         getEnv("DB_NAME", "tracker"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		IPInfoToken:    getEnv("IP_INFO_TOKEN", ""),
		JWTSecret:      getEnv("JWT_SECRET", "your_jwt_secret"),
		JWTExpiration:  getEnvAsInt("JWT_EXPIRATION", 900),
		RedisURL:       getEnv("REDIS_URL", ""),
		MailerHost:     getEnv("MAILER_HOST", ""),
		MailerPort:     getEnvAsInt("MAILER_PORT", 587),
		MailerUsername: getEnv("MAILER_USERNAME", ""),
		MailerPassword: getEnv("MAILER_PASSWORD", ""),
		EncryptionKey:  getEnv("ENCRYPTION_KEY", "your_32_byte_encryption_key"),
	}
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
