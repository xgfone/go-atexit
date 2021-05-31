// Copyright 2021 xgfone
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

// Package atexit is used to manage the exit functions of the program.
//
// Example
//
//     package main
//
//     import (
//         "flag"
//         "log"
//         "os"
//
//         "github.com/xgfone/go-atexit"
//     )
//
//     var logfile string
//
//     func init() {
//         flag.StringVar(&logfile, "logfile", "", "the log file path")
//
//         atexit.RegisterWithPriority(1, func() { log.Println("the program exits") })
//         atexit.Register(func() { log.Println("do something to clean") })
//     }
//
//     func main() {
//         flag.Parse()
//
//         if logfile != "" {
//             file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
//             if err != nil {
//                 log.Println(err)
//                 atexit.Exit(1)
//             }
//             log.SetOutput(file)
//
//             // Close the file before the program exits.
//             atexit.RegisterWithPriority(0, func() {
//                 log.Println("close the log file")
//                 file.Close()
//             })
//         }
//
//         log.Println("do jobs ...")
//
//         atexit.Exit(0) // The program exits.
//
//         // $ go run main.go
//         // 2021/05/29 08:29:14 do jobs ...
//         // 2021/05/29 08:29:14 do something to clean
//         // 2021/05/29 08:29:14 the program exits
//         //
//         // $ go run main.go -logfile test.log
//         // $ cat test.log
//         // 2021/05/29 08:29:19 do jobs ...
//         // 2021/05/29 08:29:19 do something to clean
//         // 2021/05/29 08:29:19 the program exits
//         // 2021/05/29 08:29:19 close the log file
//     }
//
package atexit

import (
	"os"
	"sort"
	"sync/atomic"
	"time"
)

var (
	atexits = make(exitFuncs, 0, 4)
	exitch  = make(chan struct{}, 1)
	exited  uint32
)

type exitFunc struct {
	Name string
	Func func()
	Prio int
}

type exitFuncs []exitFunc

func (fs exitFuncs) Len() int           { return len(fs) }
func (fs exitFuncs) Less(i, j int) bool { return fs[i].Prio < fs[j].Prio }
func (fs exitFuncs) Swap(i, j int)      { fs[i], fs[j] = fs[j], fs[i] }

var priority = int64(99)

// ExitFunc is used to customize the exit function.
var ExitFunc = os.Exit

// ExitDelay is used to wait for a delay duration before the program exits.
var ExitDelay = time.Millisecond * 100

// Wait waits until all the exit functions have finished to be executed.
func Wait() { <-exitch }

// Exit calls the exit functions in reverse and the program exits with the code.
func Exit(code int) {
	Execute()

	if ExitDelay > 0 {
		time.Sleep(ExitDelay)
	}

	ExitFunc(code)
}

// RegisterWithPriority registers the exit callback function with the priority,
// which will be called when calling Exit.
//
// Notice: the callback function with the higher priority will be executed
// preferentially.
func RegisterWithPriority(priority int, callback func()) {
	if callback == nil {
		panic("atexit.RegisterWithPriority: callback function is nil")
	}

	atexits = append(atexits, exitFunc{Prio: priority, Func: callback})
	sort.Stable(atexits)
}

// Register is the same as RegisterWithPriority, but increase the priority
// beginning with 100. For example,
//   Register(callback) // ==> RegisterWithPriority(100, callback)
//   Register(callback) // ==> RegisterWithPriority(101, callback)
func Register(callback func()) {
	RegisterWithPriority(int(atomic.AddInt64(&priority, 1)), callback)
}

// Execute calls all the registered exit functions in reverse.
//
// Notice: It only executes the exits functions once.
func Execute() {
	if atomic.CompareAndSwapUint32(&exited, 0, 1) {
		for _len := len(atexits) - 1; _len >= 0; _len-- {
			func(f func()) { defer recover(); f() }(atexits[_len].Func)
		}

		close(exitch)
	}
}
