package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect rabittmQ %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open the channel or connection")
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"email_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare the queue")
	}
	msgs, err := ch.Consume(
		queue.Name,
		"", //Consumer tag (leave blank for auto- generated)
		true,
		false, //exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register the consumer")
	}
	forever := make(chan bool)
	fmt.Println("Waiting for email messages")

	go func() {
		for msg := range msgs {
			fmt.Printf("\n Received Messages %s\n", msg.Body)
			var email EmailData
			if err := json.Unmarshal(msg.Body, &email); err != nil {
				log.Printf("Invalid email data %v", err)
				continue
			}
			sendEmail(email)
		}
	}()
	<-forever
}

func sendEmail(email EmailData) {
	m := gomail.NewMessage()
	m.SetHeader("From", "email")
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)

	m.SetBody("text/html", email.Body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "email", "password")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email")
		return
	}
	log.Printf("Email sent successfully to %s\n", email.To)

}
