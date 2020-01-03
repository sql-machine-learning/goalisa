// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goalisa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDSN(t *testing.T) {
	a := assert.New(t)
	dsn := `pid:pkey@example.com?env=sadfjkiem`
	cfg, err := ParseDSN(dsn)
	a.NoError(err)
	expected := Config{
		POPAccessID:  "pid",
		POPAccessKey: "pkey",
		POPURL:       "example.com",
		Env:          "sadfjkiem",
		Verbose:      false}
	a.Equal(expected, *cfg)
}

func TestParseDSNError(t *testing.T) {
	a := assert.New(t)
	badDSN := []string{
		`:pkey@example.com?env=sadfjkiem`,
		`pid:@example.com?env=sadfjkiem`,
		`pid:pkey@?env=sadfjkiem`,
	}
	for _, dsn := range badDSN {
		_, err := ParseDSN(dsn)
		a.Error(err)
	}
}

func TestConfig_FormatDSN(t *testing.T) {
	a := assert.New(t)
	cfg := Config{
		POPAccessID:  "pid",
		POPAccessKey: "pkey",
		POPURL:       "example.com",
		Env:          "sadfjkiem",
		Verbose:      false}
	expected := `pid:pkey@example.com?env=sadfjkiem&verbose=false`
	a.Equal(expected, cfg.FormatDSN())
}

func TestRoundTrip(t *testing.T) {
	a := assert.New(t)
	expected := `pid:pkey@example.com?env=sadfjkiem&verbose=true`
	cfg, err := ParseDSN(expected)
	a.NoError(err)
	a.Equal(expected, cfg.FormatDSN())
}
