# Golang App Shutdown Handling

This package provides convenient interface for working with `os.Signal`.

Multiple hooks can be applied, they will be called simultaneously on app shutdown.

## Sample

```go
package main

import (
	"log"
	"time"

	"github.com/ztrue/shutdown"
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
```

Find more executable examples in [examples](examples) dir.

## How to Use

Import:

```go
import "github.com/ztrue/shutdown"
```

Add shutdown hook:

```go
shutdown.Add(func() {
  log.Println("Stopping")
})
```

Remove hook:

```go
key := shutdown.Add(func() {
  log.Println("Stopping")
})

shutdown.Remove(key)
```

With custom key:

```go
shutdown.AddWithKey("mykey", func() {
  log.Println("Stopping")
})

shutdown.Remove("mykey")
```

With signal parameter:

```go
shutdown.AddWithParam(func(os.Signal) {
  log.Println("Stopping because of", os.Signal)
})
```

Listen for specific signals:

```go
shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
```
