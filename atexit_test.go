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
)

func TestFuncs(t *testing.T) {
	var skip int
	var funcs []Func
	buf := bytes.NewBuffer(nil)
	funcs = registerCallback(funcs, "test", skip, 0, func() { buf.WriteString("exit1\n") })
	funcs = registerCallback(funcs, "test", skip, 3, func() { buf.WriteString("exit2\n") })
	funcs = registerCallback(funcs, "test", skip, 3, func() { buf.WriteString("exit3\n") })
	funcs = registerCallback(funcs, "test", skip, 2, func() { buf.WriteString("exit4\n") })
	funcs = registerCallback(funcs, "test", skip, 1, func() { buf.WriteString("exit5\n") })
	funcs = registerCallback(funcs, "test", skip, 2, func() { buf.WriteString("exit6\n") })

	expectlines := []int{26, 30, 29, 31, 27, 28}
	for i, f := range funcs {
		if line := expectlines[i]; line != f.Line {
			t.Errorf("%d: expect line %d, but got %d", i, line, f.Line)
		}
	}
}
