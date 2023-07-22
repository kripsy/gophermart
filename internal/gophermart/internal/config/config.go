package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress      string
	DatabaseAddress string
	AccrualAddress  string
	LoggerLevel     string
	//<<<<<<< HEAD
	MigrationsPath string
	//=======
	PublicKey string
	//>>>>>>> dev
}

var cfg = &Config{}

var defaultPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3vebIWQERuzgn0T/70ta
9QKjuvneTt84YniM+jKtUXOjY69CrPImf6YViq+h5cjqgKKkraJD7zqRLRQH/Mlj
eZ0rSI6w9kPLtUk9BjQMyzRgaPltz4AvxaCFAA60AUhH6dpJbm1PEAFbZRwOVzwI
8faXRjbDySWb1M8gnw3aw1b2y7aSMm8OGx5+w3kxoe6L06P+b2oeZDoy8nYcf1Ef
W263+q7RsBskvbwFbCKAPT3moOV+V3Hi1Cmc+SCvHcCpvfn4UpL5nxddHJN7Ny84
F2T0uagVrGlF7BBfibtkT9RJQCq6ehr9yRA2CSZw1Fo1RdUn6SGB6CLXvNs5vQkf
TQIDAQAB
-----END PUBLIC KEY-----`

func GetConfig() *Config {
	return cfg
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

	accrualAddress := flag.String(
		"r",
		os.Getenv("ACCRUAL_SYSTEM_ADDRESS"),
		"Enter address exec http server accrual. Or use ACCRUAL_SYSTEM_ADDRESS env")

	loggerLevel := flag.String(
		"l",
		os.Getenv("LOGGER_LEVEL"),
		"Enter logger level as Warn. Or use LOGGER_LEVEL env")

	//<<<<<<< HEAD
	migrationsPath := flag.String(
		"m",
		os.Getenv("MIGRATIONS_PATH_AUTH"),
		"Enter migrations path. Or use MIGRATIONS_PATH_AUTH env")
	//=======
	publicKey := flag.String(
		"publicKey",
		os.Getenv("PUBLIC_KEY"),
		"Enter public key. Or use PUBLIC_KEY env")
	//>>>>>>> dev

	flag.Parse()

	if *runAddress == "" {
		*runAddress = "localhost:8080"
	}

	if *databaseAddress == "" {
		*databaseAddress = "postgres://gophermart:RASKkCt3PVEU@localhost:5432/auth?sslmode=disable"
	}

	if *accrualAddress == "" {
		*accrualAddress = "http://localhost:8081"
	}

	if *loggerLevel == "" {
		*loggerLevel = "Warn"
	}

	//<<<<<<< HEAD
	if *migrationsPath == "" {
		*migrationsPath = "./db/gophermart/migrations"
	}
	//=======
	if *publicKey == "" {
		*publicKey = defaultPublicKey
		//>>>>>>> dev
	}

	cfg = &Config{
		RunAddress:      *runAddress,
		DatabaseAddress: *databaseAddress,
		AccrualAddress:  *accrualAddress,
		LoggerLevel:     *loggerLevel,
		//<<<<<<< HEAD
		MigrationsPath: *migrationsPath,
		//=======
		PublicKey: *publicKey,
		//>>>>>>> dev
	}

	return cfg
}
