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

// RHeader represents Header for human readable format.
//   DiskGuid is string type.
type RHeader struct {
	Signature      uint64
	Revision       uint32
	Size           uint32
	Crc32OfHeader  uint32
	Reserved       uint32
	CurrentLBA     uint64
	BackupLBA      uint64
	FirstUsableLBA uint64
	LastUsableLBA  uint64
	DiskGuid       string
	StartingLBA    uint64
	NumOfEntries   uint32
	SizeOfEntry    uint32
	Crc32OfEntries uint32
}

func NewRHeader(h Header) *RHeader {
	ret := &RHeader{Signature: h.Signature, Revision: h.Revision, Size: h.Size, Crc32OfHeader: h.Crc32OfHeader, Reserved: h.Reserved, CurrentLBA: h.CurrentLBA, BackupLBA: h.BackupLBA, FirstUsableLBA: h.FirstUsableLBA, LastUsableLBA: h.LastUsableLBA, StartingLBA: h.StartingLBA, NumOfEntries: h.NumOfEntries, SizeOfEntry: h.SizeOfEntry, Crc32OfEntries: h.Crc32OfEntries}

	ret.DiskGuid = h.DiskGuid.String()
	return ret
}

// REntry represents Entry for human readable format.
//  TypeGuid/UniqueGuid/Name are string type.
type REntry struct {
	TypeGuid   string
	UniqueGuid string
	FirstLBA   uint64
	LastLBA    uint64
	AttrFlags  uint64
	Name       string
}

func NewREntry(e Entry) *REntry {
	ret := &REntry{FirstLBA: e.FirstLBA, LastLBA: e.LastLBA, AttrFlags: e.AttrFlags}
	ret.TypeGuid = e.TypeGuid.String()
	ret.UniqueGuid = e.UniqueGuid.String()
	ret.Name = e.ReadName()

	return ret
}

// RChs represents Chs for human readable format.
type RChs struct {
	Head     uint
	Sector   uint
	Cylinder uint
}

func NewRChs(c Chs) *RChs {
	return &RChs{Head: c.Head(), Sector: c.Sector(), Cylinder: c.Cylinder()}
}

// RMbrEntry represents MbrEntry for human readable format.
type RMbrEntry struct {
	BootFlag byte
	FirstChs RChs
	Id       byte
	LastChs  RChs
	FirstLBA uint32
	AllLBA   uint32
}

func NewRMbrEntry(m MbrEntry) *RMbrEntry {
	ret := &RMbrEntry{BootFlag: m.BootFlag, Id: m.Id, FirstLBA: m.FirstLBA, AllLBA: m.AllLBA}
	fc := NewRChs(m.FirstChs)
	ret.FirstChs = *fc

	fc = NewRChs(m.LastChs)
	ret.LastChs = *fc
	return ret
}

// RMbr represents Mbr for human readable format.
type RMbr struct {
	Entries   [4]RMbrEntry
	Signature uint16
}

func NewRMbr(m Mbr) *RMbr {
	ret := &RMbr{Signature: m.Signature}
	for i := 0; i < 4; i++ {
		e := NewRMbrEntry(m.Entries[i])
		ret.Entries[i] = *e
	}
	return ret
}

// RMbr represents RGpt for human readable format.
type RGpt struct {
	Mbr           RMbr
	Header        RHeader
	Entries       map[uint]REntry
	BackupEntries map[uint]REntry
	BackupHeader  RHeader
}

func NewRGpt(g Gpt) *RGpt {
	ret := &RGpt{}
	m := NewRMbr(g.Mbr)
	ret.Mbr = *m

	h := NewRHeader(g.Header)
	ret.Header = *h

	h = NewRHeader(g.BackupHeader)
	ret.BackupHeader = *h

	ret.Entries = make(map[uint]REntry)
	for i, v := range g.Entries {
		if !v.IsBlank() {
			e := NewREntry(g.Entries[i])
			ret.Entries[uint(i)] = *e
		}
	}

	ret.BackupEntries = make(map[uint]REntry)
	for i, v := range g.BackupEntries {
		if !v.IsBlank() {
			e := NewREntry(g.BackupEntries[i])
			ret.BackupEntries[uint(i)] = *e
		}
	}

	return ret
}
