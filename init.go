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
	initfuncs = make(priofuncs, 0, 4)
)

// RegisterInitWithPriority registers the init function with the priority,
// which will be called when calling Init.
//
// Notice: The smaller the value, the higher the priority.
func RegisterInitWithPriority(priority int, init func()) {
	initfuncs = append(initfuncs, priofunc{Prio: priority, Func: init})
	sort.Stable(initfuncs)
}

// RegisterInit is the same as RegisterInitWithPriority, but increase the priority
// starting with 100. For example,
//   Register(init) // ==> RegisterInitWithPriority(100, init)
//   Register(init) // ==> RegisterInitWithPriority(101, init)
func RegisterInit(init func()) {
	RegisterInitWithPriority(int(atomic.AddInt64(&initprio, 1)), init)
}

// OnInit is the alias of RegisterInit.
func OnInit(init func()) {
	RegisterInit(init)
}

// OnInitWithPriority is the alias of RegisterInitWithPriority.
func OnInitWithPriority(priority int, init func()) {
	RegisterInitWithPriority(priority, init)
}

// Init calls all the registered init functions.
func Init() {
	for i, _len := 0, len(initfuncs); i < _len; i++ {
		initfuncs[i].Func()
	}
}
