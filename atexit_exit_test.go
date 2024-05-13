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
	"bytes"
	"testing"
	"time"
)

func TestRegisterAndExecute(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	OnExitWithPriority(1, func() { buf.WriteString("1") })
	OnExitWithPriority(2, func() { buf.WriteString("2") })
	OnExitWithPriority(3, func() { buf.WriteString("3") })
	OnExitWithPriority(3, func() { buf.WriteString("4") })
	OnExitWithPriority(2, func() { buf.WriteString("5") })
	OnExitWithPriority(1, func() { buf.WriteString("6") })
	go func() { time.Sleep(time.Second); Execute() }()

	start := time.Now()
	Wait()
	if time.Since(start) < time.Second {
		t.Error("wait for too few seconds")
	}

	expect := "435261"
	if s := buf.String(); s != expect {
		t.Errorf("expect '%s', but got '%s'", expect, s)
	}
}
