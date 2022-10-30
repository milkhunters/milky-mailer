package main

import (
	"flag"
	"milky-mailer/internal/app"
	"milky-mailer/internal/configer"
)

func main() {
	// Init flags
	appName := flag.String("APP_NAME", "milky-mailer", "App name for consul")
	consulHost := flag.String("CONSUL_HOST", "192.168.3.41", "Consul host")
	flag.Parse()

	// Get config
	AMQPCfg, EmailCfg, err := configer.GetConfig(*appName, *consulHost)
	if err != nil {
		panic(err)
	}

	// Run app
	app.Run(AMQPCfg, EmailCfg)
}
