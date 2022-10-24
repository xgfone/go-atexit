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
	"fmt"
	"testing"
)

func testRegisterInit(f func()) {
	f()
}

func TestFuncFileLine(t *testing.T) {
	OnInit(func() {})
	OnInit(func() {})
	testRegisterInit(func() {
		OnInit(func() {})
	})

	funcs := GetAllInitFuncs()
	if funcs[0].Line != 27 {
		t.Errorf("0: expect line %d, but got %d", 27, funcs[0].Line)
	}
	if funcs[1].Line != 28 {
		t.Errorf("1: expect line %d, but got %d", 28, funcs[0].Line)
	}
	if funcs[2].Line != 30 {
		t.Errorf("2: expect line %d, but got %d", 30, funcs[0].Line)
	}
}

func ExampleInit() {
	// Register the init functions.
	OnInit(func() { fmt.Println("init1") })
	OnInit(func() { fmt.Println("init2") })
	OnInitWithPriority(20, func() { fmt.Println("init3") })
	OnInitWithPriority(10, func() { fmt.Println("init4") })
	OnInitWithPriority(10, func() { fmt.Println("init5") })
	OnInitWithPriority(30, func() { fmt.Println("init6") })
	OnInit(func() { fmt.Println("init7") })

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
