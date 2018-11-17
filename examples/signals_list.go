package main

import (
	"log"
	"syscall"
	"time"

	"github.com/ztrue/shutdown"
	// "../../shutdown"
)

func main() {
	shutdown.Add(func() {
		log.Println("Stopping...")
		time.Sleep(2 * time.Second)
		log.Println("Stopped")
	})

	// App emulation.
	go run()

	// Handle only SIGINT and SIGTERM.
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}

func run() {
	log.Println("App running, press CTRL + C to stop")
	select {}
}
