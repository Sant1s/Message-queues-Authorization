package qeue

import (
	"log"
	"strings"

	"github.com/Sant1s/MessageQueueTask/pkg/email"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func RunConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	qRegistration, err := ch.QueueDeclare(
		email.REGISTRATION_ACTON, // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgsRegistration, err := ch.Consume(
		qRegistration.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	failOnError(err, "Failed to register a consumer")

	qResetPassword, err := ch.QueueDeclare(
		email.RESET_PASSWORD_ACTION, // name
		false,                       // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgsResetPassword, err := ch.Consume(
		qResetPassword.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgsRegistration {
			err := processRegistrationMessages(string(d.Body))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		for d := range msgsResetPassword {
			err := processResetPasswordMessages(string(d.Body))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func processRegistrationMessages(msg string) error {
	data := strings.Split(msg, "&")
	login := strings.Split(data[0], "login=")
	mail := strings.Split(data[1], "email=")

	return email.SendMessageLoginUser(login[1], mail[1])
}

func processResetPasswordMessages(msg string) error {
	data := strings.Split(msg, "&")
	login := strings.Split(data[0], "login=")
	mail := strings.Split(data[1], "email=")

	return email.SendMessageResetPassoword(login[1], mail[1])
}
