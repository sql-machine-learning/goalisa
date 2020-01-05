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
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

var b64EnvStr = base64.URLEncoding.EncodeToString([]byte(`{"param1":"value1"}`))

func TestEncodeEnv(t *testing.T) {
	a := assert.New(t)
	a.Equal(b64EnvStr, encodeEnv(map[string]string{"param1": "value1"}))
}

func TestParseDSN(t *testing.T) {
	a := assert.New(t)
	dsn := `pid:pkey@example.com?env=` + b64EnvStr
	cfg, err := ParseDSN(dsn)
	a.NoError(err)
	expected := Config{
		POPAccessID:     "pid",
		POPAccessSecret: "pkey",
		POPURL:          "example.com",
		Env:             map[string]string{"param1": "value1"},
		Verbose:         false}
	a.Equal(expected, *cfg)
}

func TestParseDSNError(t *testing.T) {
	a := assert.New(t)
	badDSN := []string{
		`:pkey@example.com?env=` + b64EnvStr,
		`pid:@example.com?env=` + b64EnvStr,
		`pid:pkey@?env=` + b64EnvStr,
	}
	for _, dsn := range badDSN {
		_, err := ParseDSN(dsn)
		a.Error(err)
	}
}

func TestConfig_FormatDSN(t *testing.T) {
	a := assert.New(t)
	cfg := Config{
		POPAccessID:     "pid",
		POPAccessSecret: "pkey",
		POPURL:          "example.com",
		Env:             map[string]string{"param1": "value1"},
		Verbose:         false}
	expected := `pid:pkey@example.com?env=` + b64EnvStr + `&verbose=false`
	a.Equal(expected, cfg.FormatDSN())
}

func TestRoundTrip(t *testing.T) {
	a := assert.New(t)
	expected := `pid:pkey@example.com?env=` + b64EnvStr + `&verbose=true`
	cfg, err := ParseDSN(expected)
	a.NoError(err)
	a.Equal(expected, cfg.FormatDSN())
}
