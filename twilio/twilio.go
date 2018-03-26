package twilio

import "fmt"

// Send messages throught the Twilio service to any given mobile phone number.
func SendMessage(phoneNumber int64, message string, options ...interface{}) {
	fmt.Print(phoneNumber, message, options)
}
