package main

import (
	"log"

	"crypto/aes"
	"crypto/cipher"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var key = []byte("0123456789ABCDEF0123456789ABCDEF")

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}
func clearQueue(ch *amqp.Channel, queueName string) {
	log.Printf("Clearing the queue: %s", queueName)
	_, err := ch.QueuePurge(queueName, false)
	if err != nil {
		log.Printf("Failed to clear the queue: %s", err)
	} else {
		log.Println("Queue cleared successfully")
	}
}
func decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
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
			encryptedData := d.Body

			decryptedData, err := decrypt(encryptedData, key)
			failOnError(err, "Failed to decrypt data")

			data := string(decryptedData)

			// if data[0] == '+' {
			// 	MailBhejo.SendWhatsAppMessage(data)
			// 	err := MailBhejo.SendSms(data)
			// 	if err != nil {
			// 		fmt.Println(err)
			// 		os.Exit(1)
			// 	}

			// } else {

			MailBhejo.SendEmail(data)
			// }

			log.Printf("Done")
			d.Ack(false) // Acknowledge the message
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit, press CTRL+C")
	<-make(chan struct{})
}
