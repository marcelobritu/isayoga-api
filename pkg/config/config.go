package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	MercadoPago MercadoPagoConfig
	Telemetry   TelemetryConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	MongoURI    string
	MongoDBName string
}

type MercadoPagoConfig struct {
	AccessToken string
	NotifyURL   string
	BackURL     string
}

type TelemetryConfig struct {
	ZipkinURL      string
	ServiceName    string
	ServiceVersion string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		Database: DatabaseConfig{
			MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
			MongoDBName: getEnv("MONGO_DB_NAME", "isayoga"),
		},
		MercadoPago: MercadoPagoConfig{
			AccessToken: getEnv("MERCADOPAGO_ACCESS_TOKEN", ""),
			NotifyURL:   getEnv("MERCADOPAGO_NOTIFY_URL", ""),
			BackURL:     getEnv("MERCADOPAGO_BACK_URL", "http://localhost:8080"),
		},
		Telemetry: TelemetryConfig{
			ZipkinURL:      getEnv("ZIPKIN_URL", "http://localhost:9411/api/v2/spans"),
			ServiceName:    getEnv("SERVICE_NAME", "isayoga-api"),
			ServiceVersion: getEnv("SERVICE_VERSION", "1.0.0"),
		},
	}

	if config.Database.MongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI é obrigatório")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
