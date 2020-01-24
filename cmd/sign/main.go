package main

import (
	"fmt"
	"log"
	"meeting-app/s2t"
)

func main() {
	t, err := s2t.LoadTranscript()
	if err != nil {
		log.Fatalf("Failed to load transcript: %v", err)
		// try to create one, see speech-to-text main.go
	}
	//TODO Sign transcript
	fmt.Println(t)
}
