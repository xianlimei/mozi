// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// golang-plugins uses the new plugin feature of Go 1.8 to
// implement hot code swapping in Go.
// This is highly experimental and just a way for me to learn
// how plugins work and what limitations I find.
//
// Limitations:
//
// This only works on Linux.
// We poll regularly the plugins directory instead of using fsnotify.
// We recompile every time, even if the code has not changed.
// This causes a continuously growing memory requirement (memory leak?).
package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/chenkaiC4/golang-plugins/jober"
)

type payload struct {
	Element []int
}

func main() {
	jer := jober.NewJober(filepath.Join(".", "tasks"), filepath.Join(".", "sos"))
	jer.Start()

	// jober.
	for index := 0; index < 100; index++ {
		time.Sleep(1 * time.Second)
		fmt.Println("******* send print job *******")

		args1 := &jober.JobArgs{
			Name: "Print",
			Args: []byte("CK"),
		}
		err := jer.AddJob(args1)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(1 * time.Second)
		fmt.Println("******* send sum job *******")
		n := []int{1, 2, 3}
		pld := &payload{
			Element: n,
		}
		plgb, err := json.Marshal(pld)
		if err != nil {
			return
		}
		args2 := &jober.JobArgs{
			Name: "Sum",
			Args: plgb,
		}
		err = jer.AddJob(args2)
		if err != nil {
			fmt.Println(err)
		}
	}
	time.Sleep(30 * time.Second)

	jer.Clear()
}
