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

import "sync/atomic"

var (
	initprio  = int64(99)
	initfuncs = make([]Func, 0, 4)
)

func registerInitCallback(priority int, init func()) {
	const prefix = "atexit.OnInitWithPriority: init"
	initfuncs = registerCallback(initfuncs, prefix, 2, priority, init)
}

// GetAllInitFuncs returns all the registered init functions.
func GetAllInitFuncs() []Func {
	return append([]Func{}, initfuncs...)
}

// OnInitWithPriority registers the init function with the priority,
// which will be called when calling Init.
//
// Notice: The smaller the value, the higher the priority.
func OnInitWithPriority(priority int, init func()) {
	registerInitCallback(priority, init)
}

// OnInit is the same as OnInitWithPriority, but increase the priority
// starting with 100. For example,
//
//	OnInit(init) // ==> OnInitWithPriority(100, init)
//	OnInit(init) // ==> OnInitWithPriority(101, init)
func OnInit(init func()) {
	registerInitCallback(int(atomic.AddInt64(&initprio, 1)), init)
}

// Init calls all the registered init functions.
//
// If setting the environment variable "DEBUG" to a true bool value
// parsed by strconv.ParseBool, it will print the debug log to stdout.
func Init() { runInits(initfuncs) }
