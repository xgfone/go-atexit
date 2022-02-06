// Copyright 2021~2022 xgfone
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
	"context"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type priofunc struct {
	Func func()
	Prio int
}

type priofuncs []priofunc

func (fs priofuncs) Len() int           { return len(fs) }
func (fs priofuncs) Less(i, j int) bool { return fs[i].Prio < fs[j].Prio }
func (fs priofuncs) Swap(i, j int)      { fs[i], fs[j] = fs[j], fs[i] }

var (
	executed    uint32
	execlock    sync.Mutex
	priority    = int64(99)
	executech   = make(chan struct{})
	exitfuncs   = make(priofuncs, 0, 4)
	ctx, cancel = context.WithCancel(context.Background())
)

func execute() (yes bool) {
	if atomic.LoadUint32(&executed) == 0 {
		execlock.Lock()
		defer execlock.Unlock()
		if yes = atomic.LoadUint32(&executed) == 0; yes {
			defer atomic.StoreUint32(&executed, 1)

			cancel()
			for _len := len(exitfuncs) - 1; _len >= 0; _len-- {
				func(f func()) { defer recover(); f() }(exitfuncs[_len].Func)
			}
			close(executech)
		}
	}
	return
}

// RegisterWithPriority registers the exit callback function with the priority,
// which will be called when calling Exit.
//
// Notice: The bigger the value, the higher the priority.
func RegisterWithPriority(priority int, callback func()) {
	if callback == nil {
		panic("atexit.RegisterWithPriority: callback function is nil")
	}

	exitfuncs = append(exitfuncs, priofunc{Prio: priority, Func: callback})
	sort.Stable(exitfuncs)
}

// Register is the same as RegisterWithPriority, but increase the priority
// starting with 100. For example,
//   Register(callback) // ==> RegisterWithPriority(100, callback)
//   Register(callback) // ==> RegisterWithPriority(101, callback)
func Register(callback func()) {
	RegisterWithPriority(int(atomic.AddInt64(&priority, 1)), callback)
}

// Context returns the context to indicate whether the registered exit funtions
// are executed, that's, the function Execute is called.
func Context() context.Context { return ctx }

// Done is a convenient function that is equal to Context().Done().
//
// DEPRCATED: use Context().Done() instead.
func Done() <-chan struct{} { return Context().Done() }

// Executed reports whether the registered exit funtions have finished to execute.
func Executed() (yes bool) { return atomic.LoadUint32(&executed) == 1 }

// Execute calls all the registered exit functions in reverse.
//
// Notice: The exit functions are executed only once.
func Execute() { execute() }

// Wait waits until all the registered exit functions have finished to execute.
func Wait() { <-executech }

// ExitFunc is used to customize the exit function.
var ExitFunc = os.Exit

// ExitDelay is used to wait for a delay duration before calling ExitFunc.
var ExitDelay = time.Millisecond * 100

// Exit calls the exit functions in reverse and the program exits with the code.
func Exit(code int) {
	if executed := execute(); executed {
		if ExitDelay > 0 {
			time.Sleep(ExitDelay)
		}
		ExitFunc(code)
	}
}
