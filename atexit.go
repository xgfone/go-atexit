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
//   package main
//
//   import (
//       "flag"
//       "log"
//       "os"
//
//       "github.com/xgfone/go-atexit"
//   )
//
//   var logfile = flag.String("logfile", "", "the log file path")
//
//   func init() {
//       // Register the exit functions
//       atexit.OnExitWithPriority(1, func() { log.Println("the program exits") })
//       atexit.OnExit(func() { log.Println("do something to clean") })
//
//       // Register the init functions.
//       atexit.OnInit(flag.Parse)
//       atexit.OnInit(func() {
//           if *logfile != "" {
//               file, err := os.OpenFile(*logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
//               if err != nil {
//                   log.Println(err)
//                   atexit.Exit(1)
//               } else {
//                   log.SetOutput(file)
//               }
//
//               // Close the file before the program exits.
//               atexit.OnExitWithPriority(0, func() {
//                   log.Println("close the log file")
//                   file.Close()
//               })
//           }
//       })
//   }
//
//   func main() {
//       atexit.Init()
//
//       log.Println("do jobs ...")
//
//       atexit.Exit(0)
//
//       // $ go run main.go
//       // 2021/05/29 08:29:14 do jobs ...
//       // 2021/05/29 08:29:14 do something to clean
//       // 2021/05/29 08:29:14 the program exits
//       //
//       // $ go run main.go -logfile test.log
//       // $ cat test.log
//       // 2021/05/29 08:29:19 do jobs ...
//       // 2021/05/29 08:29:19 do something to clean
//       // 2021/05/29 08:29:19 the program exits
//       // 2021/05/29 08:29:19 close the log file
//   }
//
package atexit

import (
	"context"
	"os"
	"sort"
	"sync/atomic"
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
	priority    = int64(99)
	executech   = make(chan struct{})
	exitfuncs   = make(priofuncs, 0, 4)
	ctx, cancel = context.WithCancel(context.Background())
)

func execute() {
	if atomic.CompareAndSwapUint32(&executed, 0, 1) {
		cancel()
		for _len := len(exitfuncs) - 1; _len >= 0; _len-- {
			func(f func()) { defer recover(); f() }(exitfuncs[_len].Func)
		}
		close(executech)
	}
	return
}

// OnExitWithPriority registers the exit callback function with the priority,
// which will be called when calling Exit.
//
// Notice: The bigger the value, the higher the priority.
func OnExitWithPriority(priority int, callback func()) {
	if callback == nil {
		panic("atexit.OnExitWithPriority: callback function is nil")
	}

	exitfuncs = append(exitfuncs, priofunc{Prio: priority, Func: callback})
	sort.Stable(exitfuncs)
}

// OnExit is the same as OnExitWithPriority, but increase the priority
// starting with 100. For example,
//   OnExit(callback) // ==> OnExitWithPriority(100, callback)
//   OnExit(callback) // ==> OnExitWithPriority(101, callback)
func OnExit(callback func()) {
	OnExitWithPriority(int(atomic.AddInt64(&priority, 1)), callback)
}

// Context returns the context to indicate whether the registered exit funtions
// are executed, that's, the function Execute is called.
func Context() context.Context { return ctx }

// Done is a convenient function that is equal to Context().Done().
func Done() <-chan struct{} { return Context().Done() }

// Executed reports whether the registered exit funtions have been executed.
func Executed() bool { return atomic.LoadUint32(&executed) == 1 }

// Execute calls all the registered exit functions in reverse.
//
// Notice: The exit functions are executed only once.
func Execute() { execute() }

// Wait waits until all the registered exit functions have finished to execute.
func Wait() { <-executech }

// ExitFunc is used to customize the exit function.
var ExitFunc = os.Exit

// Exit calls the exit functions in reverse and the program exits with the code.
func Exit(code int) {
	execute()
	ExitFunc(code)
}
