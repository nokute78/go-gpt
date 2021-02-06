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
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"strings"
	"unicode/utf16"
)

// The signature of GPT Header. "EFI PART".
const HeaderSignature = 0x5452415020494645

// Header reprensents the partition table header of GPT.
// ref: https://en.wikipedia.org/wiki/GUID_Partition_Table#Partition_table_header_(LBA_1)
type Header struct {
	Signature      uint64
	Revision       uint32
	Size           uint32
	Crc32OfHeader  uint32
	Reserved       uint32
	CurrentLBA     uint64
	BackupLBA      uint64
	FirstUsableLBA uint64
	LastUsableLBA  uint64
	DiskGuid       Guid
	StartingLBA    uint64
	NumOfEntries   uint32
	SizeOfEntry    uint32
	Crc32OfEntries uint32
	Reserved2      [420]byte
}

// ReadHeader reads GPT Header from r.
// This function treats sector size is 512 byte.
// It returns GPT Header pointer or error if error occured.
func ReadHeader(r io.Reader) (*Header, error) {
	h := &Header{}

	err := binary.Read(r, binary.LittleEndian, h)
	if err != nil {
		return nil, fmt.Errorf("ReadHeader:%w", err)
	}
	return h, nil
}

// IsValid reports whether h is valid.
// It checks if
//   The signature is valid.
//   Crc32 of header is valid.
func (h Header) IsValid() bool {
	if h.Signature != HeaderSignature {
		return false
	}
	c := h.Crc32OfHeader
	h.Crc32OfHeader = 0

	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.LittleEndian, &h)
	if err != nil {
		return false
	}
	return c == crc32.ChecksumIEEE(buf.Bytes()[:h.Size])
}

// Entry represents a partition entries of GPT.
// ref: https://en.wikipedia.org/wiki/GUID_Partition_Table#Partition_entries_(LBA_2%E2%80%9333)
type Entry struct {
	TypeGuid   Guid
	UniqueGuid Guid
	FirstLBA   uint64
	LastLBA    uint64
	AttrFlags  uint64
	Name       [36]uint16
}

// ReadEntry reads GPT Entry from r.
// It returns GPT Entry pointer or error if error occured.
func ReadEntry(r io.Reader) (*Entry, error) {
	e := &Entry{}
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return nil, fmt.Errorf("ReadEntry:%w", err)
	}
	return e, nil
}

func (e Entry) IsBlank() bool {
	if e.TypeGuid.Equal(*ZeroGuid) && e.UniqueGuid.Equal(*ZeroGuid) {
		return true
	}
	return false
}

// ReadName returns partition name string.
// It converts from UTF-16 name string.
func (e Entry) ReadName() string {
	s := string(utf16.Decode(e.Name[:]))
	return strings.ReplaceAll(s, "\u0000", "") // remove null char
}

// WriteName writes s as UTF16 string.
func (e *Entry) WriteName(s string) error {
	r := utf16.Encode([]rune(s))
	if len(r) > len(e.Name) {
		return errors.New("string is too long")
	}
	for i := 0; i < len(e.Name); i++ {
		e.Name[i] = 0
	}
	for i := 0; i < len(r); i++ {
		e.Name[i] = r[i]
	}
	return nil
}

// GPT represents MBR, GPT header and each partition entries.
// It also contains backup header and entries.
type Gpt struct {
	Mbr           Mbr
	Header        Header
	Entries       []Entry
	BackupEntries []Entry
	BackupHeader  Header
}

// ReadGpt reads GPT from rs.
func ReadGpt(rs io.ReadSeeker) (*Gpt, error) {
	sectorSize := int64(512)
	g := &Gpt{}

	rs.Seek(0, io.SeekStart)
	m, err := ReadMbr(rs)
	if err != nil {
		return nil, fmt.Errorf("ReadGpt:%w", err)
	}
	g.Mbr = *m

	// Read Primary Header
	rs.Seek(sectorSize*1, io.SeekStart)

	h, err := ReadHeader(rs)
	if err != nil {
		return nil, fmt.Errorf("ReadGpt:%w", err)
	}
	g.Header = *h

	// Read Primary Entries
	rs.Seek(sectorSize*int64(g.Header.StartingLBA), io.SeekStart)
	for i := uint32(0); i < g.Header.NumOfEntries; i++ {
		e, err := ReadEntry(rs)
		if err != nil {
			return nil, fmt.Errorf("ReadGpt:%w", err)
		}
		g.Entries = append(g.Entries, *e)
	}

	// Read Backup Header
	rs.Seek(sectorSize*int64(g.Header.BackupLBA), io.SeekStart)
	h, err = ReadHeader(rs)
	if err != nil {
		return nil, fmt.Errorf("ReadGpt:%w", err)
	}
	g.BackupHeader = *h

	// Read Backup Entries
	rs.Seek(sectorSize*int64(g.BackupHeader.StartingLBA), io.SeekStart)
	for i := uint32(0); i < g.BackupHeader.NumOfEntries; i++ {
		e, err := ReadEntry(rs)
		if err != nil {
			return nil, fmt.Errorf("ReadGpt:%w", err)
		}
		g.BackupEntries = append(g.BackupEntries, *e)
	}

	return g, nil
}
