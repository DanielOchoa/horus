package twilio

import "testing"
import "os"
import "github.com/joho/godotenv"
import "github.com/DanielOchoa/horus/config"

// Notice use of .env.test file for twilio test credentials (no msn sent out)
// TODO: Change function so we pass the http client so we can stub it out.
func TestSendMessage(t *testing.T) {
	envPath := config.GetGoPath() + config.GetProjectPath() + "/.env.test"
	if err := godotenv.Load(envPath); err != nil {
		t.Error(err)
	}

	_, err := SendMessage(os.Getenv("TWILIO_TO_NUMBER"), "A test made this!")
	if err != nil {
		t.Error(err)
	}
}
