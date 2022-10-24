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

package atexit

import (
	"bytes"
	"sort"
	"testing"
	"time"
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
	if funcs[0].Line != 29 {
		t.Errorf("0: expect line %d, but got %d", 29, funcs[0].Line)
	}
	if funcs[1].Line != 30 {
		t.Errorf("1: expect line %d, but got %d", 30, funcs[0].Line)
	}
	if funcs[2].Line != 32 {
		t.Errorf("2: expect line %d, but got %d", 32, funcs[0].Line)
	}
}

func TestExitFuncs(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	exits := funcs{
		Func{Prio: 0, Func: func() { buf.WriteString("exit1\n") }},
		Func{Prio: 3, Func: func() { buf.WriteString("exit2\n") }},
		Func{Prio: 3, Func: func() { buf.WriteString("exit3\n") }},
		Func{Prio: 2, Func: func() { buf.WriteString("exit4\n") }},
		Func{Prio: 1, Func: func() { buf.WriteString("exit5\n") }},
		Func{Prio: 2, Func: func() { buf.WriteString("exit6\n") }},
	}
	sort.Stable(exits)

	for i := 0; i < 6; i++ {
		exits[i].Func()
	}

	expect := "exit1\nexit5\nexit4\nexit6\nexit2\nexit3\n"
	if result := buf.String(); result != expect {
		t.Errorf("expect '%s', but got '%s'", expect, result)
	}
}

func TestRegisterAndExecute(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	if Executed() {
		t.Errorf("expect unexecuted, but got executed")
	}

	OnExitWithPriority(1, func() { buf.WriteString("1") })
	OnExitWithPriority(2, func() { buf.WriteString("2") })
	OnExitWithPriority(3, func() { buf.WriteString("3") })
	OnExitWithPriority(3, func() { buf.WriteString("4") })
	OnExitWithPriority(2, func() { buf.WriteString("5") })
	OnExitWithPriority(1, func() { buf.WriteString("6") })
	go func() { time.Sleep(time.Second); Execute() }()

	start := time.Now()
	Wait()
	if time.Now().Sub(start) < time.Second {
		t.Error("wait for too few seconds")
	}

	if !Executed() {
		t.Errorf("expect executed, but got unexecuted")
	}

	expect := "435261"
	if s := buf.String(); s != expect {
		t.Errorf("expect '%s', but got '%s'", expect, s)
	}
}
