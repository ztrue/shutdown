package main

import (
	"log"
	"os"
	"time"

	"github.com/ztrue/shutdown"
	// "../../shutdown"
)

func main() {
	shutdown.AddWithParam(func(sig os.Signal) {
		log.Println(sig, "foo stopping...")
		log.Println(sig, "foo stopped")
	})

	bazKey := shutdown.AddWithParam(func(sig os.Signal) {
		log.Println(sig, "bar stopping...")
		time.Sleep(time.Second)
		log.Println(sig, "bar stopped")
	})

	shutdown.AddWithParam(func(sig os.Signal) {
		log.Println(sig, "baz stopping...")
		time.Sleep(2 * time.Second)
		log.Println(sig, "baz stopped")
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
