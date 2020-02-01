package main

import (
	"log"
	"meeting-app/edit"
)

func main() {
	log.Fatalf("Web server failure: %v", edit.Start())
}
