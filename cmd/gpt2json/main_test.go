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
	"bytes"
	"strings"
	"testing"
)

func TestCliRun(t *testing.T) {
	type testcase struct {
		name   string
		input  []string
		expect int
	}

	cases := []testcase{
		{"no args", []string{}, ExitArgError},
		{"show Version", []string{"-V"}, ExitOK},
		{"help", []string{"-h"}, ExitOK},
	}

	nullbuf := bytes.NewBuffer([]byte{})

	for _, v := range cases {
		args := []string{"program-name"}
		args = append(args, v.input...)

		cli := &CLI{OutStream: nullbuf, ErrStream: nullbuf, quiet: true}
		ret := cli.Run(args)
		if ret != v.expect {
			t.Errorf("%s:given %d expect %d", v.name, ret, v.expect)
		}
	}
}

func TestShowVersion(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	cli := &CLI{OutStream: buf, quiet: true}

	ret := cli.Run([]string{"showVer", "-V"})
	if ret != ExitOK {
		t.Errorf("ret is not ExitOK, ret=%d", ret)
	}

	if !strings.Contains(string(buf.Bytes()), "Ver:") {
		t.Errorf("not version string: %s", string(buf.Bytes()))
	}
}
