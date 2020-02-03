package main

import (
	"log"
	"meeting-app/sign"
)

func main() {
	log.Fatalf("Web server failure: %v", sign.Start())
}
