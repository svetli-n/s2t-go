package main

import (
	"fmt"
	"log"
	"meeting-app/platform/storage"
)

func main() {
	t, err := storage.LoadTranscript()
	if err != nil {
		log.Fatalf("Failed to load transcript: %v", err)
		// try to create one, see speech-to-text main.go
	}
	//TODO Sign transcript
	fmt.Println(t)
}
