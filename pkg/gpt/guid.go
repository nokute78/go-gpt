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
	"fmt"
)

var EspGuid *Guid

// ZeroGuid 00000000-0000-0000-0000-00000000000
var ZeroGuid *Guid

func init() {
	var err error
	EspGuid, err = NewGuidFromString("C12A7328-F81F-11D2-BA4B-00A0C93EC93B")
	if err != nil {
		panic(err)
	}
	ZeroGuid = &Guid{}
}

type Guid [16]byte

// String implements Stringer interface.
func (g Guid) String() string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%04x-%06x", g[3], g[2], g[1], g[0], g[5], g[4], g[7], g[6], g[8:10], g[10:])
}

// Equal check if gg is same guid.
func (g Guid) Equal(gg Guid) bool {
	for i, v := range g {
		if v != gg[i] {
			return false
		}
	}
	return true
}

// NewGuidFromBytes returns guid from b.
func NewGuidFromBytes(b []byte) (*Guid, error) {
	g := &Guid{}
	if len(b) != len(*g) {
		return nil, fmt.Errorf("length should be %d", len(*g))
	}
	for i, v := range b {
		g[i] = v
	}
	return g, nil
}

// NewGuidFromString returns guid grom s.
// Format is 00112233-4455-6677-8899-aabbccddeeff
func NewGuidFromString(s string) (*Guid, error) {
	// 33221100-5544-7766-8899-aabbccddeeff
	b := make([]byte, 16)
	_, err := fmt.Sscanf(s, "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x", &b[3], &b[2], &b[1], &b[0], &b[5], &b[4], &b[7], &b[6], &b[8], &b[9], &b[10], &b[11], &b[12], &b[13], &b[14], &b[15])
	if err != nil {
		return nil, fmt.Errorf("NewGuidFromString:%w", err)
	}
	return NewGuidFromBytes(b)
}
