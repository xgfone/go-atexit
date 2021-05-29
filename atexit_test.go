// Copyright 2021 xgfone
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
)

func TestExitFuncs(t *testing.T) {
	exits := exitFuncs{
		exitFunc{Name: "exit1", Prio: 0},
		exitFunc{Name: "exit2", Prio: 3},
		exitFunc{Name: "exit3", Prio: 3},
		exitFunc{Name: "exit4", Prio: 2},
		exitFunc{Name: "exit5", Prio: 1},
		exitFunc{Name: "exit6", Prio: 2},
	}
	sort.Stable(exits)

	expects := exitFuncs{
		exitFunc{Name: "exit1"},
		exitFunc{Name: "exit5"},
		exitFunc{Name: "exit4"},
		exitFunc{Name: "exit6"},
		exitFunc{Name: "exit2"},
		exitFunc{Name: "exit3"},
	}

	for i := 0; i < 6; i++ {
		if expects[i].Name != exits[i].Name {
			t.Errorf("expect '%s', but got '%s'", expects[i].Name, exits[i].Name)
		}
	}
}

func TestRegisterAndExecute(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	RegisterWithPriority(1, func() { buf.WriteString("1") })
	RegisterWithPriority(2, func() { buf.WriteString("2") })
	RegisterWithPriority(3, func() { buf.WriteString("3") })
	RegisterWithPriority(3, func() { buf.WriteString("4") })
	RegisterWithPriority(2, func() { buf.WriteString("5") })
	RegisterWithPriority(1, func() { buf.WriteString("6") })
	Execute()

	expect := "435261"
	if s := buf.String(); s != expect {
		t.Errorf("expect '%s', but got '%s'", expect, s)
	}
}