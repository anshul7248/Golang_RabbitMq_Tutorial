package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// "github.com/streadway/amqp"

func main() {
	// amqp.Dial()
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
		"email_queue", //Queue name
		true,          // Durable
		false,         // AutoDelete
		false,         // Exlusive
		false,         // NoWait
		nil,           //Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare the queue %v", err)
	}

	emailBody := `
	{
	"to":"xowop69450@limtu.com",
	"subject":"Welcome to golang and devops channel",
	"body":"Subscribe please and give motivation"
}
	`

	ch.Publish(
		"",         //Exchange name
		queue.Name, //Routing Key
		false,      //Mandatory
		false,      //Immediate
		amqp.Publishing{
			ContentType: "application/json", //tell  consumer data type
			Body:        []byte(emailBody),  //actual data in bytes
		})
	if err != nil {
		log.Fatalf("Failed to publish messages %v", err)
	}

	fmt.Println("Email request sent to queue ", emailBody)
}
