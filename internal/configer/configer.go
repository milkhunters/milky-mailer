package configer

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"milky-mailer/internal/mailer"
	"os"
	"strconv"
)

// TODO Доп переменные среды приоритетнее консула

type AMQPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Queue    string
}

func GetConfig() (*AMQPConfig, *mailer.EmailSenderConfig, error) {
	// Get data for init from flags
	appName := flag.String("APP_NAME", "", "App name for consul")
	consulAddress := flag.String("CONSUL_ADDRESS", "", "Consul address")
	flag.Parse()

	// Get data for init from env
	if os.Getenv("APP_NAME") != "" {
		*appName = os.Getenv("APP_NAME")
	}
	if os.Getenv("CONSUL_ADDRESS") != "" {
		*consulAddress = os.Getenv("CONSUL_ADDRESS")
	}

	// Set default values
	if *appName == "" {
		*appName = "milky-mailer"
	}
	if *consulAddress == "" {
		*consulAddress = "localhost:8500"
	}

	// Configure Consul connection
	consulCfg := api.DefaultConfig()
	consulCfg.Address = *consulAddress

	// Create a new client
	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, nil, err
	}

	kv := client.KV()

	// Get values from consul
	var pair *api.KVPair

	// TODO Использовать viper
	var AMQPCfg AMQPConfig
	pair, _, err = kv.Get(fmt.Sprintf("%s/amqp/host", *appName), nil)
	AMQPCfg.Host = string(pair.Value)
	pair, _, err = kv.Get(fmt.Sprintf("%s/amqp/port", *appName), nil)
	AMQPCfg.Port, err = strconv.Atoi(string(pair.Value))
	pair, _, err = kv.Get(fmt.Sprintf("%s/amqp/user", *appName), nil)
	AMQPCfg.User = string(pair.Value)
	pair, _, err = kv.Get(fmt.Sprintf("%s/amqp/password", *appName), nil)
	AMQPCfg.Password = string(pair.Value)
	pair, _, err = kv.Get(fmt.Sprintf("%s/amqp/queue", *appName), nil)
	AMQPCfg.Queue = string(pair.Value)

	var EmailCfg mailer.EmailSenderConfig
	pair, _, err = kv.Get(fmt.Sprintf("%s/email/from", *appName), nil)
	EmailCfg.From = string(pair.Value)
	pair, _, err = kv.Get(fmt.Sprintf("%s/email/user", *appName), nil)
	EmailCfg.User = string(pair.Value)
	pair, _, err = kv.Get(fmt.Sprintf("%s/email/host", *appName), nil)
	EmailCfg.Host = string(pair.Value)
	pair, _, err = kv.Get(fmt.Sprintf("%s/email/password", *appName), nil)
	EmailCfg.Password = string(pair.Value)
	pair, _, err = kv.Get(fmt.Sprintf("%s/email/port", *appName), nil)
	EmailCfg.Port, err = strconv.Atoi(string(pair.Value))

	if err != nil {
		return nil, nil, err
	}

	return &AMQPCfg, &EmailCfg, nil
}
