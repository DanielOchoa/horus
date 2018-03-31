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
// TODO: Pass in http.Client so we can stub it out for testing.
func SendMessage(phoneNumber string, message string, options ...interface{}) (map[string]interface{}, error) {

	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", phoneNumber)
	msgData.Set("From", fromNumber)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// error making request? try again? hmm
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Printf("Horus: MSN sent with id of SID: %q\n", data["sid"])
			return data, nil
		}
		return nil, err
	}
	fmt.Printf("Horus: Twilio responded with status: %q", resp.Status)
	return nil, errors.New(fmt.Sprintf("Horus: !Whops, request failed -- %q", resp.Status))
}
