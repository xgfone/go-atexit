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

// Register is the alias of OnExit.
//
// Deprecated. Please use OnExit.
func Register(callback func()) {
	OnExit(callback)
}

// RegisterWithPriority is the alias of OnExitWithPriority.
//
// Deprecated. Please use OnExitWithPriority.
func RegisterWithPriority(priority int, callback func()) {
	OnExitWithPriority(priority, callback)
}

// RegisterInitWithPriority is the alias of OnInitWithPriority.
//
// Deprecated. Please use OnInitWithPriority.
func RegisterInitWithPriority(priority int, init func()) {
	OnInitWithPriority(priority, init)
}

// RegisterInit is the alias of OnInit.
//
// Deprecated. Please use OnInit.
func RegisterInit(init func()) {
	OnInit(init)
}
