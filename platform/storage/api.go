package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const resourcesDir = "platform/storage/resources"
const transcriptFileName = "transcript.txt"
const transcriptJSONFileName = "transcript.json"
const confirmationJSONFileName = "confirmation.json"
const audioFileName = "audio-short.wav"
const emailsFileName = "emails.txt"

var audioFilePath = filepath.Join(resourcesDir, audioFileName)
var transcriptFilePath = filepath.Join(resourcesDir, transcriptFileName)
var transcriptJSONFilePath = filepath.Join(resourcesDir, transcriptJSONFileName)
var confirmationJSONFilePath = filepath.Join(resourcesDir, confirmationJSONFileName)
var emailsFilePath = filepath.Join(resourcesDir, emailsFileName)

func LoadAudio() ([]byte, error) {
	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(audioFilePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return nil, err
	}
	return data, nil
}

func SaveTranscript(transcript Transcript) error {
	f, err := os.OpenFile(transcriptJSONFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to open file: %v with err: %v", transcriptFilePath, err)
		return err
	}
	defer f.Sync()
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "\t")
	if err = e.Encode(transcript); err != nil {
		return err
	}
	return nil
}

func LoadTranscript() (Transcript, error) {
	t := Transcript{}
	trans, err := ioutil.ReadFile(transcriptJSONFilePath)
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal(trans, &t); err != nil {
		return t, err
	}
	return t, nil
}

func SaveConfirmation(confirmation Confirmation) error {
	f, err := os.OpenFile(confirmationJSONFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to open file: %v with err: %v", confirmationJSONFilePath, err)
		return err
	}
	defer f.Sync()
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "\t")
	if err = e.Encode(confirmation); err != nil {
		return err
	}
	return nil
}

func LoadConfirmation() (Confirmation, error) {
	c := Confirmation{}
	trans, err := ioutil.ReadFile(confirmationJSONFilePath)
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(trans, &c); err != nil {
		return c, err
	}
	return c, nil
}
