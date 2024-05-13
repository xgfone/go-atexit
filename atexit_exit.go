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

package atexit

import (
	"context"
	"os"
	"sync/atomic"
	"time"
)

var (
	executed    uint32
	exitprio    = int64(99)
	exitfuncs   = make([]Func, 0, 4)
	exitexecch  = make(chan struct{})
	ctx, cancel = context.WithCancel(context.Background())
)

func execute() {
	if atomic.CompareAndSwapUint32(&executed, 0, 1) {
		cancel()
		runExits(exitfuncs)
		close(exitexecch)
	}
}

func registerExitCallback(priority int, exit func()) {
	const prefix = "atexit.OnExitWithPriority: exit callback"
	exitfuncs = registerCallback(exitfuncs, prefix, 2, priority, exit)
}

// GetAllExitFuncs returns all the registered exit functions.
func GetAllExitFuncs() []Func {
	return append([]Func{}, exitfuncs...)
}

// OnExitWithPriority registers the exit callback function with the priority,
// which will be called when calling Exit.
//
// Notice: The bigger the value, the higher the priority.
func OnExitWithPriority(priority int, callback func()) {
	registerExitCallback(priority, callback)
}

// OnExit is the same as OnExitWithPriority, but increase the priority
// starting with 100. For example,
//
//	OnExit(callback) // ==> OnExitWithPriority(100, callback)
//	OnExit(callback) // ==> OnExitWithPriority(101, callback)
func OnExit(callback func()) {
	registerExitCallback(int(atomic.AddInt64(&exitprio, 1)), callback)
}

// Context returns the context to indicate whether the registered exit funtions
// are executed, that's, the function Execute is called.
func Context() context.Context { return ctx }

// Done is a convenient function that is equal to Context().Done().
func Done() <-chan struct{} { return Context().Done() }

// Execute calls all the registered exit functions in reverse.
//
// If setting the environment variable "DEBUG" to a true bool value
// parsed by strconv.ParseBool, it will print the debug log to stdout.
//
// Notice: The exit functions are executed only once.
func Execute() { execute() }

// Wait waits until all the registered exit functions have finished to execute.
func Wait() { <-exitexecch; time.Sleep(time.Millisecond * 10) }

// ExitFunc is used to customize the exit function.
var ExitFunc = os.Exit

// Exit calls the exit functions in reverse and the program exits with the code.
func Exit(code int) {
	execute()
	ExitFunc(code)
}
