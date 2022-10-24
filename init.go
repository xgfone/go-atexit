// Copyright 2022 xgfone
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
	"sort"
	"sync/atomic"
)

var (
	initprio  = int64(99)
	initfuncs = make(funcs, 0, 4)
)

func registerInitCallback(priority int, init func()) {
	if init == nil {
		panic("atexit.OnInitWithPriority: init function is nil")
	}

	file, line := getFileLine(3)
	pf := Func{Prio: priority, Func: init, Line: line, File: file}
	initfuncs = append(initfuncs, pf)
	sort.Stable(initfuncs)
}

// GetAllInitFuncs returns all the registered init functions.
func GetAllInitFuncs() []Func {
	funcs := make(funcs, len(initfuncs))
	for i, f := range initfuncs {
		funcs[i] = f
	}
	return funcs
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
//	Register(init) // ==> OnInitWithPriority(100, init)
//	Register(init) // ==> OnInitWithPriority(101, init)
func OnInit(init func()) {
	registerInitCallback(int(atomic.AddInt64(&initprio, 1)), init)
}

// Init calls all the registered init functions.
func Init() {
	for i, _len := 0, len(initfuncs); i < _len; i++ {
		initfuncs[i].Func()
	}
}
