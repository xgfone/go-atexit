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

import "fmt"

func ExampleInit() {
	// Register the init functions.
	RegisterInit(func() { fmt.Println("init1") })
	RegisterInit(func() { fmt.Println("init2") })
	RegisterInitWithPriority(20, func() { fmt.Println("init3") })
	RegisterInitWithPriority(10, func() { fmt.Println("init4") })
	RegisterInitWithPriority(10, func() { fmt.Println("init5") })
	RegisterInitWithPriority(30, func() { fmt.Println("init6") })
	RegisterInit(func() { fmt.Println("init7") })

	// Call the registered init functions.
	Init()

	// Output:
	// init4
	// init5
	// init3
	// init6
	// init1
	// init2
	// init7
}
