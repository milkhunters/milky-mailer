package main

import (
	"milky-mailer/internal/app"
	"milky-mailer/internal/configer"
	"os"
	"strconv"
)

func main() {

	consulHost := os.Getenv("CONSUL_HOST")
	consulRoot := os.Getenv("CONSUL_ROOT")
	consulPort, err := strconv.Atoi(os.Getenv("CONSUL_PORT"))
	if err != nil {
		panic(err)
	}

	config, err := configer.NewConfig(consulHost, consulRoot, consulPort)
	if err != nil {
		panic(err)
	}

	err = app.Run(config)
	if err != nil {
		panic(err)
	}

}
