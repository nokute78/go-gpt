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
	"strings"
	"testing"
)

const testdir = "testdata"

func ReadHeaderData(t *testing.T) (*gpt.Header, error) {
	t.Helper()
	f, err := os.Open(filepath.Join(testdir, "header.bin"))
	if err != nil {
		t.Fatalf("os.Open err:%s", err)
	}
	defer f.Close()
	return gpt.ReadHeader(f)
}

func TestReadHeader(t *testing.T) {
	h, err := ReadHeaderData(t)
	if err != nil {
		t.Fatalf("ReadHeaderData err:%s", err)
	}
	if h.Signature != gpt.HeaderSignature {
		t.Errorf("signature is invalid: given :0x%x\n expect:0x%x", h.Signature, gpt.HeaderSignature)
	}
}

func TestHeaderIsValid(t *testing.T) {
	h, err := ReadHeaderData(t)
	if err != nil {
		t.Fatalf("ReadHeader err:%s", err)
	}
	if !h.IsValid() {
		t.Errorf("It should be valid")
	}

	h.Signature = 0
	if h.IsValid() {
		t.Errorf("true? The signature is invalid")
	}
	h.Signature = gpt.HeaderSignature

	h.Revision = 0
	if h.IsValid() {
		t.Errorf("true? Revision is invalid")
	}
}

func ReadEntryData(t *testing.T, name string) (*gpt.Entry, error) {
	t.Helper()
	f, err := os.Open(filepath.Join(testdir, name))
	if err != nil {
		t.Fatalf("os.Open err:%s", err)
	}
	defer f.Close()
	return gpt.ReadEntry(f)
}

func TestReadEntry(t *testing.T) {
	_, err := ReadEntryData(t, "esp_entry.bin")
	if err != nil {
		t.Fatalf("ReadEntryData err:%s", err)
	}
}

func TestReadEntryName(t *testing.T) {
	e, err := ReadEntryData(t, "esp_entry.bin")
	if err != nil {
		t.Fatalf("ReadEntryData err:%s", err)
	}

	expect := "EFI System"
	if !strings.Contains(e.ReadName(), expect) {
		t.Errorf("Name mismatch\n given :\"%s\"\n expect:\"%s\"", e.ReadName(), expect)
	}
}

func TestWriteEntryName(t *testing.T) {
	e, err := ReadEntryData(t, "esp_entry.bin")
	if err != nil {
		t.Fatalf("ReadEntryData err:%s", err)
	}

	name := "Sample Partition"
	if err := e.WriteName(name); err != nil {
		t.Errorf("WriteName error: %s", err)
	}
	if !strings.Contains(e.ReadName(), name) {
		t.Errorf("Name mismatch\n given :\"%s\"\n expect:\"%s\"", e.ReadName(), name)
	}

	too_long_name := strings.Repeat("a", 128)
	if err := e.WriteName(too_long_name); err == nil {
		t.Errorf("Too Long Name:it should be error\n")
	}
}

func TestReadGpt(t *testing.T) {
	f, err := os.Open(filepath.Join(testdir, "gpt_sample.bin"))
	if err != nil {
		t.Fatalf("os.Open err:%s", err)
	}
	defer f.Close()
	_, err = gpt.ReadGpt(f)
	if err != nil {
		t.Fatalf("ReadGpt err:%s", err)
	}
}
