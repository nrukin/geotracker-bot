package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("token not set")
	}
	token := os.Args[1]
	a, err := NewApp(token, "track.db", true)
	if err != nil {
		log.Fatal(err)
	}
	a.Start()
}
