package app

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"milky-mailer/internal/configer"
	"milky-mailer/internal/mailer"
	"os"
	"strings"
)

func Run(config *configer.Config) error {

	conn, err := amqp.DialConfig(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d",
			config.AMQP.User,
			config.AMQP.Password,
			config.AMQP.Host,
			config.AMQP.Port),
		amqp.Config{
			Vhost: config.AMQP.VHost,
		})
	if err != nil {
		return errors.Join(err, errors.New("failed to connect to RabbitMQ"))
	}
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	if err != nil {
		return errors.Join(err, errors.New("failed to open a channel"))
	}
	defer amqpChannel.Close()

	_, err = amqpChannel.QueueDeclare(
		config.AMQP.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Join(err, errors.New("failed to declare a queue"))
	}

	// create exchange
	err = amqpChannel.ExchangeDeclare(
		config.AMQP.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Join(err, errors.New("failed to declare a exchange"))
	}

	// bind queue to exchange
	err = amqpChannel.QueueBind(
		config.AMQP.Queue,
		"",
		config.AMQP.Exchange,
		false,
		nil,
	)
	if err != nil {
		return errors.Join(err, errors.New("failed to bind a queue"))
	}

	err = amqpChannel.Qos(1, 0, false)
	if err != nil {
		return errors.Join(err, errors.New("failed to set QoS"))
	}

	messageChannel, err := amqpChannel.Consume(
		config.AMQP.Queue,
		"milky-mailer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Join(err, errors.New("failed to register a consumer"))
	}

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for message := range messageChannel {

			// Verify message
			var err error

			if message.Headers["To"] == nil {
				err = errors.Join(err, errors.New("header 'To' is empty"))
			}
			if message.Headers["Subject"] == nil {
				err = errors.Join(err, errors.New("header 'Subject' is empty"))
			}
			if message.Headers["FromId"] == nil {
				err = errors.Join(err, errors.New("header 'FromId' is empty"))
			}
			if message.ContentType != "text/plain" && message.ContentType != "text/html" {
				err = errors.Join(err, errors.New("content type is not supported"))
			}
			if message.Body == nil {
				err = errors.Join(err, errors.New("body is empty"))
			}

			// Check that FromId exist in config
			if _, ok := config.Senders[message.Headers["FromId"].(string)]; !ok {
				err = errors.Join(err, errors.New(fmt.Sprintf("sender '%s' is not exist in config", message.Headers["FromId"].(string))))
			}

			if err != nil {
				log.Printf("Error verefy message: %s", err)

				if err := message.Reject(false); err != nil {
					log.Printf("Error rejecting message : %s", err)
				} else {
					log.Printf("Rejected message")
				}
			}

			// Log all information about message
			log.Printf(fmt.Sprintf(
				"Received a message: \n "+
					"To: %s \n "+
					"Subject: %s \n "+
					"ContentType: %s \n "+
					"FromId: %s \n "+
					"MessageId: %s \n "+
					"Timestamp: %s \n "+
					"AppId: %s \n "+
					"Body: %s",

				message.Headers["To"].(string),
				message.Headers["Subject"].(string),
				message.ContentType,
				message.Headers["FromId"].(string),
				message.MessageId,
				message.Timestamp.String(),
				message.AppId,
				string(message.Body),
			))

			// Send email
			err = mailer.Send(
				config.Senders[message.Headers["FromId"].(string)],
				message.Headers["To"].(string),
				message.Headers["Subject"].(string),
				message.ContentType,
				string(message.Body),
			)
			if err != nil {
				log.Printf("Error send mail: %s", err)

				fmt.Println(err.Error())
				if strings.Contains(err.Error(), "550") {
					fmt.Println("Message will be deleted from queue")
					// reject without requeue
					if err := message.Reject(false); err != nil {
						log.Printf("Error rejecting message : %s", err)
					} else {
						log.Printf("Rejected message")
					}

					continue
				} else {
					// Reject message and requeue
					if err := message.Reject(true); err != nil {
						log.Printf("Error rejecting message : %s", err)
					} else {
						log.Printf("Rejected message")
					}
					continue
				}
			}

			fmt.Println("Email send")

			// Acknowledge that email was send
			if err := message.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	// Stop for program termination
	<-stopChan
	return nil
}
