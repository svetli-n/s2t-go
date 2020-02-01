package main

import (
	"log"
	"meeting-app/s2t"
)

func main() {
	log.Fatalf("Web server failure: %v", s2t.Start())
}
