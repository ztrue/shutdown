package main

import (
	"log"
	"time"

	"github.com/ztrue/shutdown"
	// "../../shutdown"
)

func main() {
	shutdown.Add(func() {
		// Write log.
		// Stop writing files.
		// Close connections.
		// Etc.
		log.Println("Stopping...")
		log.Println("3")
		time.Sleep(time.Second)
		log.Println("2")
		time.Sleep(time.Second)
		log.Println("1")
		time.Sleep(time.Second)
		log.Println("0, stopped")
	})

	// App emulation.
	go func() {
		log.Println("App running, press CTRL + C to stop")
		select {}
	}()

	shutdown.Listen()
}
