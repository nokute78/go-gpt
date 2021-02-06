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

// JsonHeader represents Header for json.
//   DiskGuid is string type.
type JsonHeader struct {
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

func NewJsonHeader(h Header) *JsonHeader {
	ret := &JsonHeader{Signature: h.Signature, Revision: h.Revision, Size: h.Size, Crc32OfHeader: h.Crc32OfHeader, Reserved: h.Reserved, CurrentLBA: h.CurrentLBA, BackupLBA: h.BackupLBA, FirstUsableLBA: h.FirstUsableLBA, LastUsableLBA: h.LastUsableLBA, StartingLBA: h.StartingLBA, NumOfEntries: h.NumOfEntries, SizeOfEntry: h.SizeOfEntry, Crc32OfEntries: h.Crc32OfEntries}

	ret.DiskGuid = h.DiskGuid.String()
	return ret
}

// JsonEntry represents Entry for json.
//  TypeGuid/UniqueGuid/Name are string type.
type JsonEntry struct {
	TypeGuid   string
	UniqueGuid string
	FirstLBA   uint64
	LastLBA    uint64
	AttrFlags  uint64
	Name       string
}

func NewJsonEntry(e Entry) *JsonEntry {
	ret := &JsonEntry{FirstLBA: e.FirstLBA, LastLBA: e.LastLBA, AttrFlags: e.AttrFlags}
	ret.TypeGuid = e.TypeGuid.String()
	ret.UniqueGuid = e.UniqueGuid.String()
	ret.Name = e.ReadName()

	return ret
}

// JsonChs represents Chs for json.
type JsonChs struct {
	Head     uint
	Sector   uint
	Cylinder uint
}

func NewJsonChs(c Chs) *JsonChs {
	return &JsonChs{Head: c.Head(), Sector: c.Sector(), Cylinder: c.Cylinder()}
}

// JsonMbrEntry represents MbrEntry for json.
type JsonMbrEntry struct {
	BootFlag byte
	FirstChs JsonChs
	Id       byte
	LastChs  JsonChs
	FirstLBA uint32
	AllLBA   uint32
}

func NewJsonMbrEntry(m MbrEntry) *JsonMbrEntry {
	ret := &JsonMbrEntry{BootFlag: m.BootFlag, Id: m.Id, FirstLBA: m.FirstLBA, AllLBA: m.AllLBA}
	fc := NewJsonChs(m.FirstChs)
	ret.FirstChs = *fc

	fc = NewJsonChs(m.LastChs)
	ret.LastChs = *fc
	return ret
}

// JsonMbr represents Mbr for json.
type JsonMbr struct {
	Entries   [4]JsonMbrEntry
	Signature uint16
}

func NewJsonMbr(m Mbr) *JsonMbr {
	ret := &JsonMbr{Signature: m.Signature}
	for i := 0; i < 4; i++ {
		e := NewJsonMbrEntry(m.Entries[i])
		ret.Entries[i] = *e
	}
	return ret
}

// JsonMbr represents JsonGpt for json.
type JsonGpt struct {
	Mbr           JsonMbr
	Header        JsonHeader
	Entries       map[uint]JsonEntry
	BackupEntries map[uint]JsonEntry
	BackupHeader  JsonHeader
}

func NewJsonGpt(g Gpt) *JsonGpt {
	ret := &JsonGpt{}
	m := NewJsonMbr(g.Mbr)
	ret.Mbr = *m

	h := NewJsonHeader(g.Header)
	ret.Header = *h

	h = NewJsonHeader(g.BackupHeader)
	ret.BackupHeader = *h

	ret.Entries = make(map[uint]JsonEntry)
	for i, v := range g.Entries {
		if !v.IsBlank() {
			e := NewJsonEntry(g.Entries[i])
			ret.Entries[uint(i)] = *e
		}
	}

	ret.BackupEntries = make(map[uint]JsonEntry)
	for i, v := range g.BackupEntries {
		if !v.IsBlank() {
			e := NewJsonEntry(g.BackupEntries[i])
			ret.BackupEntries[uint(i)] = *e
		}
	}

	return ret
}
