package edit

import (
	"errors"
	"fmt"
	"log"
	"meeting-app/platform/storage"
	"net/http"
	"strings"
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
		log.Println(t)
		emailPrticipants(t)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func emailPrticipants(t storage.Transcript) error {
	errs := make(map[string]string)
	for _, email := range t.Emails {
		if err := sendSignEmail(email, t.Edited); err != nil {
			errs[email] = err.Error()
		}
	}
	if len(errs) > 0 {
		errStr := make([]string, len(errs))
		for k, val := range errs {
			errStr = append(errStr, k, ":", val, " ")
		}
		return errors.New(strings.Join(errStr, " "))
	}
	return nil
}

func sendSignEmail(email, transcript string) error {
	return nil
}
