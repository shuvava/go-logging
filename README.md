# Golang logging package

[![Go Reference](https://pkg.go.dev/badge/github.com/shuvava/go-logging.svg)](https://pkg.go.dev/github.com/shuvava/go-logging)
![Build Status](https://github.com/shuvava/go-logging/actions/workflows/makefile.yml/badge.svg)

This logging package is a wrapper around the standard library's log package, giving you abstractions for logging at different levels.

## Installation

```shell
  go get github.com/shuvava/go-logging/logging@v1.0.0
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
