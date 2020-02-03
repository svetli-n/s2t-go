package edit

import (
	"encoding/base64"
	"fmt"
	"log"
	"meeting-app/platform/storage"
	"net/http"
	"net/smtp"
	"os"
)

const (
	smtpAddr     = "smtp.gmail.com"
	smtpAddrPort = smtpAddr + ":587"
)

func Start() error {
	http.HandleFunc("/edit", edit)

	return http.ListenAndServe(":8788", nil)
}

func edit(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		t, err := storage.LoadTranscript()
		if err != nil {
			log.Fatalf("Failed to load transcript: %v", err)
		}
		//w.Write([]byte(t.Raw))
		fmt.Fprintf(w, "<html><body><form action=\"/edit\" method=\"POST\">"+
			"<textarea rows=\"10\" cols=\"50\" name=\"edited\">%s</textarea><br>"+
			"<input type=\"submit\" value=\"Save\">"+
			"</form></body></html>", t.Raw)
	} else if req.Method == "POST" {
		t, err := storage.LoadTranscript()
		if err != nil {
			log.Fatalf("Failed to load transcript: %v", err)
		}
		t.Edited = req.FormValue("edited")
		if err = storage.SaveTranscript(t); err != nil {
			log.Fatalf("Failed to save transcript: %v", err)
		}
		if err = emailParticipants(t); err != nil {
			log.Println(err.Error())
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func emailParticipants(t storage.Transcript) error {
	a := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PWD"), smtpAddr)
	for _, email := range t.Emails {
		msg := t.Edited + "\n" + "Confirm: http://localhost:8789/sign?email=" + base64.StdEncoding.EncodeToString([]byte(email))
		if err := smtp.SendMail(smtpAddrPort, a, "meetingapp3@gmail.com", t.Emails, []byte(msg)); err != nil {
			return err
		}
	}
	return nil
}
