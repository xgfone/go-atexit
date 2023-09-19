# Go AtExit [![Build Status](https://github.com/xgfone/go-atexit/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-atexit/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-atexit)](https://pkg.go.dev/github.com/xgfone/go-atexit) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-atexit/master/LICENSE)

The package `atexit` is used to manage the exit functions of the program. Support `Go1.8+`.

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

var logfile = flag.String("logfile", "", "the log file path")

func init() {
	// Register the exit functions
	atexit.OnExitWithPriority(1, func() { log.Println("the program exits") })
	atexit.OnExit(func() { log.Println("do something to clean") })

	// Register the init functions.
	atexit.OnInit(flag.Parse)
	atexit.OnInit(func() {
		if *logfile != "" {
			file, err := os.OpenFile(*logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Println(err)
				atexit.Exit(1)
			} else {
				log.SetOutput(file)
			}

			// Close the file before the program exits.
			atexit.OnExitWithPriority(0, func() {
				log.Println("close the log file")
				file.Close()
			})
		}
	})
}

func main() {
	atexit.Init()

	log.Println("do jobs ...")

	atexit.Exit(0)

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
