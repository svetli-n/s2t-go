package s2t

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

const resourcesDir = "s2t/resources"
const transcriptFileName = "transcript.txt"
const audioFileName = "audio-short.wav"
const emailsFileName = "emails.txt"

var audioFilePath = filepath.Join(resourcesDir, audioFileName)
var transcriptFilePath = filepath.Join(resourcesDir, transcriptFileName)
var emailsFilePath = filepath.Join(resourcesDir, emailsFileName)

func Convert(audio []byte) ([]string, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16,
			LanguageCode: "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audio},
		},
	})
	if err != nil {
		log.Fatalf("Failed to recognize: %v", err)
		return nil, err
	}

	var trans []string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			trans = append(trans, alt.Transcript)
			if err != nil {
				log.Fatalf("Failed to write string: %v with error: %v", alt.Transcript, err)
				return nil, err
			}
		}
	}

	return trans, nil
}
func LoadAudio() ([]byte, error) {
	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(audioFilePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return nil, err
	}
	return data, nil
}

func WriteTranscript(transcript []string, emails []string) error {
	f, err := os.OpenFile(transcriptFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to open file: %v with err: %v", transcriptFilePath, err)
		return err
	}
	defer f.Sync()
	defer f.Close()

	for _, line := range transcript {
		_, err := f.WriteString(line)
		if err != nil {
			return err
		}
	}

	f, err = os.OpenFile(emailsFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to open file: %v with err: %v", transcriptFilePath, err)
		return err
	}
	defer f.Sync()
	defer f.Close()

	for _, email := range emails {
		_, err := f.WriteString(email + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadTranscript() (string, error) {
	t, err := ioutil.ReadFile(transcriptFilePath)
	if err != nil {
		return "", err
	}
	return string(t), nil
}
