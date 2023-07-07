package config

import (
	"flag"
	"os"
)

type Сonfiguration struct {
	FlagRunAddr     string
	FlagDataBaseDSN string
}

var config = &Сonfiguration{}

func GetConfig() *Сonfiguration {
	return config
}

func ParseFlags() {
	flag.StringVar(&config.FlagRunAddr, "a", os.Getenv("RUN_ADDRESS"), "address and port to run server")
	flag.StringVar(&config.FlagDataBaseDSN, "d", os.Getenv("DATABASE_URI"), "postgres database DNS")

	flag.Parse()
}
