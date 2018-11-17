package main

import (
	"log"
	"time"

	"github.com/ztrue/shutdown"
	// "../../shutdown"
)

func main() {
	shutdown.Add(func() {
		log.Println("foo stopping...")
		log.Println("foo stopped")
	})

	bazKey := shutdown.Add(func() {
		log.Println("bar stopping...")
		time.Sleep(time.Second)
		log.Println("bar stopped")
	})

	shutdown.Add(func() {
		log.Println("baz stopping...")
		time.Sleep(2 * time.Second)
		log.Println("baz stopped")
	})

	// App emulation.
	go run()

	shutdown.Remove(bazKey)

	shutdown.Listen()
}

func run() {
	log.Println("App running, press CTRL + C to stop")
	select {}
}
