// Copyright 2023 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package atexit is used to manage the init and exit functions of the program.
//
// Example
//
//	package main
//
//	import (
//	    "flag"
//	    "log"
//	    "os"
//
//	    "github.com/xgfone/go-atexit"
//	)
//
//	var logfile = flag.String("logfile", "", "the log file path")
//
//	func init() {
//	    // Register the exit functions
//	    atexit.OnExitWithPriority(1, func() { log.Println("the program exits") })
//	    atexit.OnExit(func() { log.Println("do something to clean") })
//
//	    // Register the init functions.
//	    atexit.OnInit(flag.Parse)
//	    atexit.OnInit(func() {
//	        if *logfile != "" {
//	            file, err := os.OpenFile(*logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
//	            if err != nil {
//	                log.Println(err)
//	                atexit.Exit(1)
//	            } else {
//	                log.SetOutput(file)
//	            }
//
//	            // Close the file before the program exits.
//	            atexit.OnExitWithPriority(0, func() {
//	                log.Println("close the log file")
//	                file.Close()
//	            })
//	        }
//	    })
//	}
//
//	func main() {
//	    atexit.Init()
//
//	    log.Println("do jobs ...")
//
//	    atexit.Exit(0)
//
//	    // $ go run main.go
//	    // 2021/05/29 08:29:14 do jobs ...
//	    // 2021/05/29 08:29:14 do something to clean
//	    // 2021/05/29 08:29:14 the program exits
//	    //
//	    // $ go run main.go -logfile test.log
//	    // $ cat test.log
//	    // 2021/05/29 08:29:19 do jobs ...
//	    // 2021/05/29 08:29:19 do something to clean
//	    // 2021/05/29 08:29:19 the program exits
//	    // 2021/05/29 08:29:19 close the log file
//	}
package atexit

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var debug bool

func init() { debug, _ = strconv.ParseBool(os.Getenv("DEBUG")) }

// SetDebug sets the debug mode.
//
// Default: parse env var "DEBUG" as bool.
func SetDebug(b bool) { debug = b }

// Func represents an init or exit function.
type Func struct {
	Func func()
	File string
	Line int
	Prio int
}

func (f Func) runInit() { f.print("init"); f.Func() }
func (f Func) runExit() { defer f.wrapPanic(); f.print("exit"); f.Func() }
func (f Func) wrapPanic() {
	if r := recover(); r != nil {
		const msg = "exit func panics: file=%s, line=%d, panic=%v\n"
		fmt.Fprintf(os.Stderr, msg, f.File, f.Line, r)
	}
}

func (f Func) print(ftype string) {
	if debug {
		fmt.Printf("run %s func: file=%s, line=%d\n", ftype, f.File, f.Line)
	}
}

func sortfuncs(funcs []Func) {
	sort.SliceStable(funcs, func(i, j int) bool {
		return funcs[i].Prio < funcs[j].Prio
	})
}

func runInits(funcs []Func) {
	for i, _len := 0, len(funcs); i < _len; i++ {
		funcs[i].runInit()
	}
}

func runExits(funcs []Func) {
	for _len := len(funcs) - 1; _len >= 0; _len-- {
		funcs[_len].runExit()
	}
}

func registerCallback(funcs []Func, prefix string, skip, priority int, f func()) []Func {
	if f == nil {
		panic(prefix + " function is nil")
	}

	file, line := getFileLine(skip + 2)
	funcs = append(funcs, Func{Prio: priority, Func: f, Line: line, File: file})
	sortfuncs(funcs)
	return funcs
}

var trimPrefixes = []string{"/pkg/mod/", "/src/"}

func getFileLine(skip int) (file string, line int) {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		for _, mark := range trimPrefixes {
			if index := strings.Index(file, mark); index > -1 {
				file = file[index+len(mark):]
				break
			}
		}
	} else {
		file = "??"
	}

	return
}
