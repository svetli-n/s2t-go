package s2t

import (
	"cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	speech2 "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io/ioutil"
	"log"
	"meeting-app/platform/storage"
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
	orgEmail := req.Form.Get("org_email")
	if err != nil {
		log.Fatalf("Failed to receive audio: %v", err)
	}
	log.Printf("File name: %v", header.Filename)
	log.Printf("File size: %v", header.Size)
	log.Printf("Emails: %v", emails)
	log.Printf("Organizer Email: %v", orgEmail)

	audio, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read audio: %v", err)
	}

	t, err := Convert(audio)
	if err != nil {
		log.Fatalf("Failed to convert audio: %v", err)
	}

	//TODO handle unique
	emailArr := strings.Split(emails, ",")

	transcript := storage.Transcript{Raw: t, Edited: "", Emails: emailArr, OrgEmail: orgEmail}
	if err = storage.SaveTranscript(transcript); err != nil {
		log.Fatalf("Failed to save to json: %v", err)
	}
	confirmation := storage.Confirmation{Emails: make(map[string]bool), Total: len(emailArr), Done: false}
	if err = storage.SaveConfirmation(confirmation); err != nil {
		log.Fatalf("Failed to save to json: %v", err)
	}
	w.WriteHeader(http.StatusOK)

}

func Convert(audio []byte) (string, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return "", err
	}

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speech2.RecognizeRequest{
		Config: &speech2.RecognitionConfig{
			Encoding:     speech2.RecognitionConfig_LINEAR16,
			LanguageCode: "en-US",
		},
		Audio: &speech2.RecognitionAudio{
			AudioSource: &speech2.RecognitionAudio_Content{Content: audio},
		},
	})
	if err != nil {
		log.Fatalf("Failed to recognize: %v", err)
		return "", err
	}

	var trans []string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			trans = append(trans, alt.Transcript)
			if err != nil {
				log.Fatalf("Failed to write string: %v with error: %v", alt.Transcript, err)
				return "", err
			}
		}
	}

	return strings.Join(trans, " "), nil
}
