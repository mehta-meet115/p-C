package Service

import (
	"context"
	"log"

	"github.com/Rhymen/go-whatsapp"
)

func sendWhatsAppMessage(phoneNumber string, message string) {
	wac, err := whatsapp.NewConn(5 * 60 * 1000)
	if err != nil {
		log.Fatal(err)
	}

	qr := make(chan string)
	go func() {
		<-qr
	}()

	session, err := wac.Login(qr)
	if err != nil {
		log.Fatal(err)
	}

	contact := phoneNumber + "@s.whatsapp.net"
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJID: contact,
		},
		Text: message,
	}

	err = wac.Send(context.Background(), text)
	if err != nil {
		log.Fatal(err)
	}

	session, err = wac.Logout()
	if err != nil {
		log.Fatal(err)
	}
}

func WhatsAppProducer(phoneNumber string) {
	message := "This is a sample WhatsApp message"
	sendWhatsAppMessage(phoneNumber, message)
}
