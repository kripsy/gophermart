package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress      string
	DatabaseAddress string
	LoggerLevel     string
}

func InitConfig() *Config {
	runAddress := flag.String(
		"a",
		"localhost:8080",
		"Enter bind address for serve auth server as ip_address:port. Or use RUN_ADDRESS env")

	databaseAddress := flag.String(
		"d",
		"postgres://gophermart:RASKkCt3PVEU@localhost:5432/auth?sslmode=disable",
		"Enter address exec http server as postgres://username:password@hostname:portNumber/databaseName?sslmode=disable. Or use DATABASE_URI env")

	loggerLevel := flag.String(
		"l",
		"Warn",
		"Enter logger level as Warn. Or use LOGGER_LEVEL env")

	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		*runAddress = envRunAddress
	}

	if envDatabaseAddress := os.Getenv("DATABASE_URI"); envDatabaseAddress != "" {
		*databaseAddress = envDatabaseAddress
	}

	if envLoggerLevel := os.Getenv("LOGGER_LEVEL"); envLoggerLevel != "" {
		*loggerLevel = envLoggerLevel
	}

	return &Config{
		RunAddress:      *runAddress,
		DatabaseAddress: *databaseAddress,
		LoggerLevel:     *loggerLevel,
	}
}
