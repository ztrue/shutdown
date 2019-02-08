# Golang App Shutdown Hooks

[![GoDoc](https://godoc.org/github.com/ztrue/shutdown?status.svg)](https://godoc.org/github.com/ztrue/shutdown)
[![Report](https://goreportcard.com/badge/github.com/ztrue/shutdown)](https://goreportcard.com/report/github.com/ztrue/shutdown)
[![Coverage Status](https://coveralls.io/repos/github/ztrue/shutdown/badge.svg?branch=master)](https://coveralls.io/github/ztrue/shutdown?branch=master)
[![Build Status](https://travis-ci.com/ztrue/shutdown.svg?branch=master)](https://travis-ci.com/ztrue/shutdown)

This package provides convenient interface for working with `os.Signal`.

Multiple hooks can be applied, they will be called simultaneously on app shutdown.

## Example

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

### Import

```go
import "github.com/ztrue/shutdown"
```

### Add Shutdown Hook

```go
shutdown.Add(func() {
	log.Println("Stopping")
})
```

### Remove Hook

```go
key := shutdown.Add(func() {
	log.Println("Stopping")
})

shutdown.Remove(key)
```

### Hook With Custom Key

```go
shutdown.AddWithKey("mykey", func() {
	log.Println("Stopping")
})

shutdown.Remove("mykey")
```

### Hook With Signal Parameter

```go
shutdown.AddWithParam(func(os.Signal) {
	log.Println("Stopping because of", os.Signal)
})
```

### Listen for Specific Signals

```go
shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
```
