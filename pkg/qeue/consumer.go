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
	conn, err := amqp.Dial(HOST_NAME)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	qRegistration, err := ch.QueueDeclare(
		email.REGISTRATION_ACTON,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgsRegistration, err := ch.Consume(
		qRegistration.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	qResetPassword, err := ch.QueueDeclare(
		email.RESET_PASSWORD_ACTION,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgsResetPassword, err := ch.Consume(
		qResetPassword.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	qPassword, err := ch.QueueDeclare(
		"password",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgsPassword, err := ch.Consume(
		qPassword.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
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

	go func() {
		for d := range msgsPassword {
			err := processResetPassword(string(d.Body))
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

	return email.SendMessageResetPassword(login[1], mail[1])
}

func processResetPassword(msg string) error {
	data := strings.Split(msg, "&")
	login := strings.Split(data[0], "login=")
	mail := strings.Split(data[1], "email=")
	password := strings.Split(data[2], "password=")

	return email.SetPassword(login[1], mail[1], password[1])
}
