package main

import (
	"flag"
	"milky-mailer/internal/app"
	"milky-mailer/internal/configer"
)

func main() {
	// Init flags
	appName := flag.String("APP_NAME", "milky-mailer", "App name for consul")
	consulAddress := flag.String("CONSUL_ADDRESS", "localhost:8500", "Consul address")
	flag.Parse()

	// Get config
	AMQPCfg, EmailCfg, err := configer.GetConfig(*appName, *consulAddress)
	if err != nil {
		panic(err)
	}

	// Run app
	app.Run(AMQPCfg, EmailCfg)
}
