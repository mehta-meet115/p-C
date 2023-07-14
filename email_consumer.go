package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	MailBhejo "teamSJ/Services"
)

var key = []byte("0123456789ABCDEF0123456789ABCDEF")

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}

func decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decryptedData := make([]byte, len(ciphertext))
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(decryptedData, ciphertext)

	return decryptedData, nil
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			encryptedData, err := hex.DecodeString(string(d.Body))
			if err != nil {
				log.Println("Failed to decode the encrypted data")
				d.Ack(false)
				continue
			}

			decryptedData, err := decrypt(encryptedData, key)
			if err != nil {
				log.Println("Failed to decrypt the data")
				d.Ack(false)
				continue
			}

			emailAddress := string(decryptedData)
			MailBhejo.SendEmail(emailAddress)

			log.Printf("Done")
			d.Ack(false) // Acknowledge the message
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit, press CTRL+C")
	<-make(chan struct{})
}
