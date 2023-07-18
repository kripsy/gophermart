package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	RunAddress      string
	DatabaseAddress string
	LoggerLevel     string
	MigrationsPath  string
	PrivateKey      string
	PublicKey       string
	TokenExp        time.Duration
}

func InitConfig() *Config {
	runAddress := flag.String(
		"a",
		"localhost:8080",
		"Enter bind address for serve auth server as ip_address:port. Or use RUN_ADDRESS env")

	databaseAddress := flag.String(
		"d",
		"postgres://gophermart:RASKkCt3PVEU@localhost:5432/gophermart?sslmode=disable",
		"Enter address exec http server as postgres://username:password@hostname:portNumber/databaseName?sslmode=disable. Or use DATABASE_URI env")

	loggerLevel := flag.String(
		"l",
		"Warn",
		"Enter logger level as Warn. Or use LOGGER_LEVEL env")

	migrationsPath := flag.String(
		"m",
		"./db/auth/migrations",
		"Enter migrations path. Or use MIGRATIONS_PATH_AUTH env")

	privateKey := flag.String(
		"privateKey",
		`
		-----BEGIN PRIVATE KEY-----
		MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDe95shZARG7OCf
		RP/vS1r1AqO6+d5O3zhieIz6Mq1Rc6Njr0Ks8iZ/phWKr6HlyOqAoqStokPvOpEt
		FAf8yWN5nStIjrD2Q8u1ST0GNAzLNGBo+W3PgC/FoIUADrQBSEfp2klubU8QAVtl
		HA5XPAjx9pdGNsPJJZvUzyCfDdrDVvbLtpIybw4bHn7DeTGh7ovTo/5vah5kOjLy
		dhx/UR9bbrf6rtGwGyS9vAVsIoA9Peag5X5XceLUKZz5IK8dwKm9+fhSkvmfF10c
		k3s3LzgXZPS5qBWsaUXsEF+Ju2RP1ElAKrp6Gv3JEDYJJnDUWjVF1SfpIYHoIte8
		2zm9CR9NAgMBAAECggEBAI3+EYUKNM8WO1YykurJintN2wdP6QtBjJ7pNp5/d3DP
		u9XX3xZUf7/6/Oz9PJUhhnW1HjqVg73uBlY2039goUDpno7ukDPEqQ4iPgKdUyh1
		ipBPiGcEs2ef+hM3SdsnNOTwZqM0aY0/z/xsCZX0XZ359Ax7A+QtVzgHUDb6k76h
		irJaTMmrtTzmTCnm72tOGhX91QusLcefffZToPjEPRlNazxeH/wdkNbYMc1GmNVH
		D/Sq3esk22t/cpeImKtv7LhShd0NCbPtM8lIJY++cBOmsM5UaQRaXrc5NSV1o1sZ
		Zjr+xJN0//p0TJAB06qhb7XsSCy4zHPXZC0cNleWiwECgYEA+GFvnP3ux6QL32Jh
		zrf6cz1P277BH/NWXKZpuqMnUiBfsHETjOASIPGytc6RL1jbTn0thkg3tt1SG1ja
		9k5pHHCtQxY5MiRpXNT52Bbu1Ko7i+e9rG7mfaws2ItWRtQSYArWj8pTMj0NpzD+
		kesgDBc67U/Gl4skLSm2WbNS9TUCgYEA5c6Ya31E9LXP4WCuBjTNDVMzb5+G7iBV
		TXQH69GhWaLTuYE1B1BR87lVj/ZVXRY9wmpQRHlmABfZPEvA77HladKUUA4oX35E
		vzjBvBp8WnuYk1VLR6vl+9AO1GLGnOoi9sYazl8l8KW+5WDHCEyQN0TnrG+jrfpw
		HelbtvFXvLkCgYEA2RNnFcEEuCySR8hXDPDUDXVvXvEHHmJwfwbd7sT675blqnIZ
		EQ0gKvSyKJ0BXGz/NkjGyc5CCyrAwK/Wpl9/E+ESPEim8kDKaNymAwp/7xNceXiu
		144RGZKpmxOj8sET0iaGwSKltYmQbieuxV7GImsHEDKhsP5lPqdu/FRyU2UCgYEA
		rEodf8jlH8oHVmNTVRfU+756+57QXEslaPIq1iPOIhOvRI6YISmYp281tL7r9OQt
		3Uozb4LMdBltJoVs2se2xYW45+QVZLKX+/0jUlFRFc0/8IWr8MnxnL65v4VmflIT
		cIvJoRs4qJi66+GIlrJAFQ+12VPBlTgDQomn1xpNuxECgYEAvQRbKiKSRRSoph7n
		Nm4DHwrquwarjwcynPPiufl3vZiJpgd6Zn6/cRpS4J6JUqSVTO0O1c9G6/EM2zsk
		tRgiBEdw38SbenkxGkTZt4kNV95qzowO2Svd+5l2nM3mfBMAAHcToqtvq0WQNNZ1
		67yGLJcuzkkgH15moAdLrN2qNQU=
		-----END PRIVATE KEY-----
		`,
		"Enter private key. Or use PRIVATE_KEY env")

	publicKey := flag.String(
		"publicKey",
		`-----BEGIN PUBLIC KEY-----
		MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3vebIWQERuzgn0T/70ta
		9QKjuvneTt84YniM+jKtUXOjY69CrPImf6YViq+h5cjqgKKkraJD7zqRLRQH/Mlj
		eZ0rSI6w9kPLtUk9BjQMyzRgaPltz4AvxaCFAA60AUhH6dpJbm1PEAFbZRwOVzwI
		8faXRjbDySWb1M8gnw3aw1b2y7aSMm8OGx5+w3kxoe6L06P+b2oeZDoy8nYcf1Ef
		W263+q7RsBskvbwFbCKAPT3moOV+V3Hi1Cmc+SCvHcCpvfn4UpL5nxddHJN7Ny84
		F2T0uagVrGlF7BBfibtkT9RJQCq6ehr9yRA2CSZw1Fo1RdUn6SGB6CLXvNs5vQkf
		TQIDAQAB
		-----END PUBLIC KEY-----
		`,
		"Enter public key. Or use PUBLIC_KEY env")

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

	if envMigrationsPath := os.Getenv("MIGRATIONS_PATH_AUTH"); envMigrationsPath != "" {
		*migrationsPath = envMigrationsPath
	}

	if envPrivateKey := os.Getenv("PRIVATE_KEY"); envPrivateKey != "" {
		*privateKey = envPrivateKey
	}

	if envPublicKey := os.Getenv("PUBLIC_KEY"); envPublicKey != "" {
		*publicKey = envPublicKey
	}

	return &Config{
		RunAddress:      *runAddress,
		DatabaseAddress: *databaseAddress,
		LoggerLevel:     *loggerLevel,
		MigrationsPath:  *migrationsPath,
		TokenExp:        time.Hour * 24,
		PrivateKey:      *privateKey,
		PublicKey:       *publicKey,
	}
}
