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

func TestJsonEntryName(t *testing.T) {
	e, err := ReadEntryData(t, "esp_entry.bin")
	if err != nil {
		t.Fatalf("ReadEntryData err:%s", err)
	}

	j := gpt.NewJsonEntry(*e)
	if !strings.Contains(j.Name, "EFI System") {
		t.Errorf("j.Name mismatch\n given :%s expect:%s", j.Name, "EFI System")
	}
}
