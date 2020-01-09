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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var b64EnvStr = base64.RawURLEncoding.EncodeToString([]byte(`{"param1":"value1"}`))

func TestEncodeEnv(t *testing.T) {
	a := assert.New(t)
	a.Equal(b64EnvStr, encodeJSONB64(map[string]string{"param1": "value1"}))
}

func TestParseDSN(t *testing.T) {
	a := assert.New(t)
	dsn := `pid:pkey@example.com?curr_project=proj&scheme=http&env=` + b64EnvStr + "&with=" + b64EnvStr
	cfg, err := ParseDSN(dsn)
	a.NoError(err)
	expected := Config{
		POPAccessID:     "pid",
		POPAccessSecret: "pkey",
		POPURL:          "example.com",
		POPScheme:       "http",
		Env:             map[string]string{"param1": "value1"},
		With:            map[string]string{"param1": "value1"},
		Verbose:         false,
		Project:         "proj"}
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

func TestDSNConfig(t *testing.T) {
	a := assert.New(t)
	dsn := "pid:psc@dw.a.hk?scheme=http&verbose=true&curr_project=jtest_env&env=eyJTS1lORVRfT05EVVRZIjoiU0tZIiwiU0tZTkVUX0FDQ0VTU0lEIjoiU0tZIiwiU0tZTkVUX1NZU1RFTUlEIjoiU0tZIiwiQUxJU0FfVEFTS19JRCI6IkFMSSIsIlNLWU5FVF9FTkRQT0lOVCI6IlNLWSIsIlNLWU5FVF9TWVNURU1fRU5WIjoiU0tZIiwiU0tZTkVUX0JJWkRBVEUiOiJTS1kiLCJTS1lORVRfQUNDRVNTS0VZIjoiU0tZIiwiU0tZTkVUX1BBQ0tBR0VJRCI6IlNLWSIsIkFMSVNBX1RBU0tfRVhFQ19UQVJHRVQiOiJBTEkifQ&with=eyJFeGVjIjoid2VjLnNoIiwiUGx1Z2luTmFtZSI6IndwZSIsIkN1c3RvbWVySWQiOiJ3Y2QifQ"
	cfg, err := ParseDSN(dsn)
	a.NoError(err)
	a.Equal("jtest_env", cfg.Project)
	a.Equal("http", cfg.POPScheme)
	a.Equal("wcd", cfg.With["CustomerId"])
	a.Equal("wpe", cfg.With["PluginName"])
	a.Equal("wec.sh", cfg.With["Exec"])
	a.Equal("SKY", cfg.Env["SKYNET_ACCESSKEY"])
}

func TestConfig_FormatDSN(t *testing.T) {
	a := assert.New(t)
	cfg := Config{
		POPAccessID:     "pid",
		POPAccessSecret: "pkey",
		POPURL:          "example.com",
		POPScheme:       "http",
		Env:             map[string]string{"param1": "value1"},
		With:            map[string]string{"param1": "value1"},
		Verbose:         false,
		Project:         "proj"}
	expected := `pid:pkey@example.com?env=` + b64EnvStr + "&with=" + b64EnvStr + `&verbose=false&curr_project=proj&scheme=http`
	a.Equal(expected, cfg.FormatDSN())
}

func TestRoundTrip(t *testing.T) {
	a := assert.New(t)
	expected := `pid:pkey@example.com?env=` + b64EnvStr + `&with=` + b64EnvStr + `&verbose=true&curr_project=proj&scheme=http`
	cfg, err := ParseDSN(expected)
	a.NoError(err)
	a.Equal(expected, cfg.FormatDSN())
}

func newConfigFromEnv(t *testing.T) *Config {
	if os.Getenv("POP_SECRET") == "" {
		t.Skip()
	}
	popURL := os.Getenv("POP_URL")
	popID := os.Getenv("POP_ID")
	popSecret := os.Getenv("POP_SECRET")
	popScheme := "http"
	verbose := os.Getenv("VERBOSE") == "true"
	envs := map[string]string{
		"SKYNET_ONDUTY":          os.Getenv("SKYNET_ONDUTY"),
		"SKYNET_ACCESSID":        os.Getenv("SKYNET_ACCESSID"),
		"SKYNET_ACCESSKEY":       os.Getenv("SKYNET_ACCESSKEY"),
		"SKYNET_ENDPOINT":        os.Getenv("SKYNET_ENDPOINT"),
		"SKYNET_SYSTEMID":        os.Getenv("SKYNET_SYSTEMID"),
		"SKYNET_PACKAGEID":       os.Getenv("SKYNET_PACKAGEID"),
		"SKYNET_SYSTEM_ENV":      os.Getenv("SKYNET_SYSTEM_ENV"),
		"SKYNET_BIZDATE":         os.Getenv("SKYNET_BIZDATE"),
		"ALISA_TASK_EXEC_TARGET": os.Getenv("ALISA_TASK_EXEC_TARGET"),
	}
	with := map[string]string{
		"CustomerId": os.Getenv("CustomerId"),
		"PluginName": os.Getenv("PluginName"),
		"Exec":       os.Getenv("Exec"),
	}
	proj := envs["SKYNET_PACKAGEID"]
	if len(envs["SKYNET_SYSTEMID"]) > 0 {
		proj += "_" + envs["SKYNET_SYSTEMID"]
	}
	return &Config{POPAccessID: popID, POPAccessSecret: popSecret, POPURL: popURL, POPScheme: popScheme, Verbose: verbose, Env: envs, Project: proj, With: with}
}
