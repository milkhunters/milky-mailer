package main

import (
	"milky-mailer/internal/app"
	"milky-mailer/internal/configer"
)

func main() {
	// Get config
	AMQPCfg, EmailCfg, err := configer.GetConfig()
	if err != nil {
		panic(err)
	}

	// Run app
	app.Run(AMQPCfg, EmailCfg)
}
