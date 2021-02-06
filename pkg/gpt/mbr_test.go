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
	"os"
	"path/filepath"
	"testing"
)

func readMbr(t *testing.T) (*gpt.Mbr, error) {
	t.Helper()
	f, err := os.Open(filepath.Join(testdir, "mbr.bin"))
	if err != nil {
		t.Fatalf("os.Open err:%s", err)
	}
	defer f.Close()
	return gpt.ReadMbr(f)
}

func TestReadMbr(t *testing.T) {
	_, err := readMbr(t)
	if err != nil {
		t.Fatalf("readMbr: %s", err)
	}
}

func TestNewChs(t *testing.T) {
	c, err := gpt.NewChs(0x100, 0, 0)
	if err == nil {
		t.Errorf("It should be error. Head > 0xff.")
	}
	c, err = gpt.NewChs(0, 0x40, 0)
	if err == nil {
		t.Errorf("It should be error. Sector > 0x3f.")
	}
	c, err = gpt.NewChs(0, 0, 0x400)
	if err == nil {
		t.Errorf("It should be error. Cylinder > 0x3ff.")
	}

	head := uint(0x3f)
	sector := uint(0x33)
	cylinder := uint(0x200)

	c, err = gpt.NewChs(head, sector, cylinder)
	if err != nil {
		t.Fatalf("NewChs err:%s", err)
	}
	if c.Head() != head {
		t.Errorf("head mismatch.\n given :0x%x\n expect:0x%x", c.Head(), head)
	}
	if c.Sector() != sector {
		t.Errorf("sector mismatch.\n given :0x%x\n expect:0x%x", c.Sector(), sector)
	}
	if c.Cylinder() != cylinder {
		t.Errorf("cylinder mismatch.\n given :0x%x\n expect:0x%x", c.Cylinder(), cylinder)
	}

}

func TestMbrIsValid(t *testing.T) {
	m := &gpt.Mbr{}
	if m.IsValid() {
		t.Errorf("It should be invalid. signature is 0")
	}
}
