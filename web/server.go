package web

import (
	"io/ioutil"
	"log"
	"meeting-app/s2t"
	"net/http"
	"strings"
)

func Start() error {
	http.HandleFunc("/upload", upload)

	return http.ListenAndServe(":8787", nil)
}

func upload(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(10_000)

	f, header, err := req.FormFile("audio")
	emails := req.Form.Get("emails")
	if err != nil {
		log.Fatalf("Failed to receive audio: %v", err)
	}
	log.Printf("File name: %v", header.Filename)
	log.Printf("File size: %v", header.Size)
	log.Printf("File size: %v", emails)

	audio, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read audio: %v", err)
	}

	t, err := s2t.Convert(audio)
	if err != nil {
		log.Fatalf("Failed to convert audio: %v", err)
	}

	err = s2t.WriteTranscript(t, strings.Split(emails, ","))
	if err != nil {
		log.Fatalf("Failed to write transcript: %v", err)
	}

	w.WriteHeader(http.StatusOK)

}
