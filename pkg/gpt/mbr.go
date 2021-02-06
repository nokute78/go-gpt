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

package gpt

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Chs represents Cylinder head sector
// ref: https://en.wikipedia.org/wiki/Cylinder-head-sector
type Chs struct {
	Body [3]byte
}

// NewChs returns Chs from head, sector and cylinder value.
func NewChs(h uint, s uint, c uint) (*Chs, error) {
	ret := &Chs{}
	if h > 0xff {
		return nil, fmt.Errorf("h = 0x%x > 0xff", h)
	}
	if s > 0x3f {
		return nil, fmt.Errorf("s = 0x%x > 0x3f", s)
	}
	if c > 0x3ff {
		return nil, fmt.Errorf("c = 0x%x > 0x3ff", c)
	}

	ret.Body[0] = byte(h)
	ret.Body[1] = byte(s) | byte((c&0x300)>>2)
	ret.Body[2] = byte(c & 0xff)

	return ret, nil
}

// Head returns the head of disk.
func (c Chs) Head() uint {
	return uint(c.Body[0])
}

// Sector returns the sector of disk.
func (c Chs) Sector() uint {
	return uint(c.Body[1] & 0x3f)
}

// Cylinder returns the sector of disk
func (c Chs) Cylinder() uint {
	return uint(c.Body[2]) | (uint(c.Body[1])&0xc0)<<2
}

// String implements fmt.Stringer interface
func (c Chs) String() string {
	return fmt.Sprintf("{cylinder:0x%x head:0x%x sector:0x%x}", c.Cylinder(), c.Head(), c.Sector())
}

// MbrEntry represents the partition entry.
type MbrEntry struct {
	BootFlag byte
	FirstChs Chs
	Id       byte
	LastChs  Chs
	FirstLBA uint32
	AllLBA   uint32
}

func (m MbrEntry) IdString() string {
	switch m.Id {
	case 0x82:
		return "Linux Swap"
	case 0x83:
		return "Linux"
	case 0xee:
		return "GPT"
	case 0xef:
		return "ESP"
	}
	return "Unknown"
}

// Mbr represents entier MBR.
// refs: https://en.wikipedia.org/wiki/Master_boot_record
type Mbr struct {
	BootCode  [446]byte
	Entries   [4]MbrEntry
	Signature uint16
}

// ReadMbr reads MBR from r.
func ReadMbr(r io.Reader) (*Mbr, error) {
	m := &Mbr{}
	err := binary.Read(r, binary.LittleEndian, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// IsValid check if the Mbr is valid.
//  The signature is 0xaa55
func (m Mbr) IsValid() bool {
	if m.Signature != 0xaa55 {
		return false
	}
	return true
}
