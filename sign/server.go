package sign

import (
	"encoding/base64"
	"log"
	"meeting-app/platform/storage"
	"net/http"
	"net/smtp"
	"os"
)

var confs = make(map[string]int)

func init() {
	t, err := storage.LoadTranscript()
	if err != nil {
		log.Fatalf("Failed to load transcript: %v", err)
	}
	for _, email := range t.Emails {
		confs[email] = 0
	}
}

func Start() error {
	http.HandleFunc("/sign", sign)
	return http.ListenAndServe(":8789", nil)
}

func sign(w http.ResponseWriter, req *http.Request) {
	emailBase64, ok := req.URL.Query()["email"]
	if !ok {
		log.Fatalf("Failed to receive email, not provided")
	}
	email, err := base64.StdEncoding.DecodeString(emailBase64[0])
	if err != nil {
		log.Fatalf("Failed to decode email: %v", err)
	}
	c, err := storage.LoadConfirmation()
	if err != nil {
		log.Fatalf("Failed to load confimration: %v", err)
	}

	c.Emails[string(email)] = true

	t, err := storage.LoadTranscript()
	if err != nil {
		log.Fatalf("Failed to load transcript: %v", err)
	}

	if !c.Done && len(c.Emails) == c.Total {
		a := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PWD"), "smtp.gmail.com")
		msg := "All confirmed."
		if err := smtp.SendMail("smtp.gmail.com:587", a, "meetingapp3@gmail.com", t.Emails, []byte(msg)); err != nil {
		}
	}

	c.Done = true
	storage.SaveConfirmation(c)
}
