# Go AtExit [![Build Status](https://github.com/xgfone/go-atexit/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-atexit/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-atexit)](https://pkg.go.dev/github.com/xgfone/go-atexit) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-atexit/master/LICENSE)

The package `atexit` is used to manage the exit functions of the program. Support `Go1.7+`.

## Install
```shell
$ go get -u github.com/xgfone/go-atexit
```

## Example
```go
package main

import (
	"flag"
	"log"
	"os"

	"github.com/xgfone/go-atexit"
)

var logfile string

func init() {
	flag.StringVar(&logfile, "logfile", "", "the log file path")

	atexit.RegisterWithPriority(1, func() { log.Println("the program exits") })
	atexit.Register(func() { log.Println("do something to clean") })
}

func main() {
	flag.Parse()

	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Println(err)
			atexit.Exit(1)
		}
		log.SetOutput(file)

		// Close the file before the program exits.
		atexit.RegisterWithPriority(0, func() {
			log.Println("close the log file")
			file.Close()
		})
	}

	log.Println("do jobs ...")

	atexit.Exit(0) // The program exits.

	// $ go run main.go
	// 2021/05/29 08:29:14 do jobs ...
	// 2021/05/29 08:29:14 do something to clean
	// 2021/05/29 08:29:14 the program exits
	//
	// $ go run main.go -logfile test.log
	// $ cat test.log
	// 2021/05/29 08:29:19 do jobs ...
	// 2021/05/29 08:29:19 do something to clean
	// 2021/05/29 08:29:19 the program exits
	// 2021/05/29 08:29:19 close the log file
}
```
