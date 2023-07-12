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
		os.Getenv("RUN_ADDRESS"),
		"Enter bind address for serve auth server as ip_address:port. Or use RUN_ADDRESS env")

	databaseAddress := flag.String(
		"d",
		os.Getenv("DATABASE_URI"),
		"Enter address exec http server as postgres://username:password@hostname:portNumber/databaseName?sslmode=disable. Or use DATABASE_URI env")

	loggerLevel := flag.String(
		"l",
		os.Getenv("LOGGER_LEVEL"),
		"Enter logger level as Warn. Or use LOGGER_LEVEL env")

	flag.Parse()

	if *runAddress == "" {
		*runAddress = "localhost:8080"
	}

	if *databaseAddress == "" {
		*databaseAddress = "postgres://gophermart:RASKkCt3PVEU@localhost:5432/auth?sslmode=disable"
	}

	if *loggerLevel == "" {
		*loggerLevel = "Warn"
	}

	return &Config{
		RunAddress:      *runAddress,
		DatabaseAddress: *databaseAddress,
		LoggerLevel:     *loggerLevel,
	}
}