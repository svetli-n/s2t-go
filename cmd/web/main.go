package main

import (
	"log"
	"meeting-app/web"
)

func main() {

	//(base) ➜  resources git:(master) ✗ ls
	//audio-short.wav
	//(base) ➜  resources git:(master) ✗ curl -F 'audio=@audio-short.wav' -F 'emails=a@a.com,b@b.com' localhost:8787/upload
	//(base) ➜  resources git:(master) ✗ ls
	//audio-short.wav emails.txt      transcript.txt
	//(base) ➜  resources git:(master) ✗

	log.Fatalf("Web server failure: %v", web.Start())
}
