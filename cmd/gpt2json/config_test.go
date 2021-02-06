/*
   Copyright 2021 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"flag"
	"testing"
)

func TestConfigure(t *testing.T) {
	type testcase struct {
		name   string
		input  []string
		expect error
	}

	cases := []testcase{
		{"no args", []string{}, ConfigNoArgs},
		{"help", []string{"-h"}, flag.ErrHelp},
		{"version", []string{"-V"}, nil},
		{"unknown opt", []string{"unknown"}, nil},
	}

	for _, v := range cases {
		_, err := Configure(v.input, true)
		if err != v.expect {
			t.Errorf("%s:given %s expect %s", v.name, err, v.expect)
		}
	}
}
