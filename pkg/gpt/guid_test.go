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

package gpt_test

import (
	"github.com/nokute78/go-gpt/pkg/gpt"
	"strings"
	"testing"
)

func TestNewGuidFromBytes(t *testing.T) {
	b := make([]byte, 16)

	g, err := gpt.NewGuidFromBytes(b)
	if err != nil {
		t.Fatalf("NewGuidFromBytes err:%s", err)
	}

	expect := "00000000-0000-0000-0000-000000000000"
	if strings.Compare(g.String(), expect) != 0 {
		t.Errorf("string mismatch:\n given :%s\n expect:%s", g, expect)
	}
}

func TestNewGuidFromString(t *testing.T) {
	g, err := gpt.NewGuidFromString("C12A7328-F81F-11D2-BA4B-00A0C93EC93B")
	if err != nil {
		t.Fatalf("NewGuidFromString err:%s", err)
	}

	expect := gpt.Guid{0x28, 0x73, 0x2a, 0xc1, 0x1f, 0xf8, 0xd2, 0x11, 0xba, 0x4b, 0x00, 0xa0, 0xc9, 0x3e, 0xc9, 0x3b}
	if !g.Equal(expect) {
		t.Errorf("string mismatch:\n given :%s\n expect:%s", g, expect)
	}
}
