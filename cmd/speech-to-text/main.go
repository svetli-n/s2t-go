package main

import (
	"log"
	"meeting-app/s2t"
)

func main() {
	audio, err := s2t.LoadAudio()
	if err != nil {
		log.Fatalf("Failed to load audio file: %v", err)
	}
	transcript, err := s2t.Convert(audio)
	if err != nil {
		log.Fatalf("Failed to convert audio to text: %v", err)
	}
	if err := s2t.WriteTranscript(transcript); err != nil {
		log.Fatalf("Failed to write transcript to file : %v", err)
	}
}
