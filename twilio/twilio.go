package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Send messages throught the Twilio service to any given mobile phone number.
func SendMessage(phoneNumber int64, message string, options ...interface{}) (map[string]interface{}, error) {

	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", fmt.Sprintf("+1%d", phoneNumber))
	msgData.Set("From", fromNumber)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
			return data, nil
		}
		return nil, err
	}
	fmt.Printf("Horus: %q", resp.Status)
	return nil, errors.New(fmt.Sprintf("Horus: !Whops, request failed -- %q", resp.Status))
}
