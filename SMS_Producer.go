package main

import (
	"log"
	"net/http"
	"net/url"
	"sync"
)

var queue = make(chan string)
var wg sync.WaitGroup

func sendSMS(phoneNumber string, message string) {
	accountSid := "your-twilio-account-sid"
	authToken := "your-twilio-auth-token"
	sender := "your-twilio-phone-number"

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", phoneNumber)
	msgData.Set("From", sender)
	msgData.Set("Body", message)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, nil)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = msgData

	_, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

func SMSProducer(phoneNumber string) {
	defer wg.Done()
	message := "This is a sample SMS"
	sendSMS(phoneNumber, message)
}

func main() {
	// Example usage of the SMSProducer function
	phoneNumbers := []string{
		"+1234567890",
		"+9876543210",
		"+1111111111",
		// Add more phone numbers to the list
	}

	// Start a goroutine for each phone number in the queue
	for _, phoneNumber := range phoneNumbers {
		wg.Add(1)
		go SMSProducer(phoneNumber)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
