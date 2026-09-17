package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string
    JWTSecret  string
    Port       string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️  No .env file found, using environment variables")
    }
    
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBUser:     getEnv("DB_USER", "root"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "atk_enterprise"),
        DBPort:     getEnv("DB_PORT", "3306"),
        JWTSecret:  getEnv("JWT_SECRET", "default-secret-key"),
        Port:       getEnv("PORT", "8080"),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}