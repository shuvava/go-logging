# Golang logging package
This logging package is a wrapper around the standard library's log package, giving you abstractions for logging at different levels.

## Installation

```shell
  go get github.com/shuvava/go-logging/logging
```

## Usage

```go
package main

import (
	"github.com/shuvava/go-logging/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logger.NewLogrusLogger(logrus.DebugLevel)
	log.Trace("Something very low level.")
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
	// Calls os.Exit(1) after logging
	log.Fatal("Bye.")
}
```
