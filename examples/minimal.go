package main

import (
	"log"

	"github.com/ztrue/shutdown"
	// "../../shutdown"
)

func main() {
	shutdown.Add(func() {
		log.Println("Stopped")
	})

	// App emulation.
	go func() {
		log.Println("App running, press CTRL + C to stop")
		select {}
	}()

	shutdown.Listen()
}
