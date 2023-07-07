package main

import (
	"github.com/kripsy/gophermart/cmd/gophermart/config"
)

func main() {
	config.ParseFlags()
	cfg := config.GetConfig()

}
