package app

// TODO заменить amqp на amqp091-go
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"milky-mailer/internal/configer"
	"milky-mailer/internal/mailer"
	"os"
)

// TODO Плохое решение. Надо переделать
func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Run(cfg *configer.AMQPConfig, mailerCfg *mailer.EmailSenderConfig) {

	// TODO Экранизировать строки
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port))
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not declare queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {

			// TODO Использовать нативный ContentType
			// Prepare message data
			mailData := mailer.EmailData{
				To:          d.Headers["To"].(string),
				Subject:     d.Headers["Subject"].(string),
				ContentType: d.Headers["ContentType"].(string),
				FromName:    d.Headers["FromName"].(string),
				Body:        string(d.Body),
			}

			// Log message
			log.Printf(fmt.Sprintf("Received a message: \n To: %s \n Subject: %s \n ContentType: %s \n FromName: %s",
				mailData.To,
				mailData.Subject,
				mailData.ContentType,
				mailData.FromName,
			))

			// Send email
			err = mailer.SendEmail(mailerCfg, &mailData)
			if err != nil {
				log.Printf("Error send mail: %s", err)
				return
			}

			// Acknowledge that email was send
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	// Stop for program termination
	<-stopChan
}
