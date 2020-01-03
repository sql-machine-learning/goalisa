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
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

var (
	reDSN = regexp.MustCompile(`^([a-zA-Z0-9_-]+):([=a-zA-Z0-9_-]+)@([:a-zA-Z0-9/_.-]+)\?([^/]+)$`)
)

// Config is the deserialization of connect string, the connection string should of format:
// pop_access_id:pop_access_key@pop_url?env=..
type Config struct {
	// POP config
	POPAccessID  string
	POPAccessKey string
	POPURL       string
	// Environment variable JSON encoded in base64 format.
	// This variable should be passed through to the http request
	Env string
	// verbose denotes whether to print logs to the terminal
	Verbose bool
}

// ParseDSN deserialize the connect string
func ParseDSN(dsn string) (*Config, error) {
	sub := reDSN.FindStringSubmatch(dsn)
	if len(sub) != 5 {
		return nil, fmt.Errorf(`dsn %s doesn't match pop_access_id:pop_access_key@pop_url?params`, dsn)
	}
	pid, pkey, purl := sub[1], sub[2], sub[3]

	kvs, err := url.ParseQuery(sub[4])
	if err != nil {
		return nil, err
	}

	requiredParameter := []string{"env"}
	for _, k := range requiredParameter {
		v := kvs.Get(k)
		if v == "" {
			return nil, fmt.Errorf(`dsn is missing required parameter %s`, k)
		}
	}

	verbose := false
	if kvs.Get("verbose") == "true" {
		verbose = true
	}

	return &Config{
		POPAccessID: pid, POPAccessKey: pkey, POPURL: purl,
		Env:     kvs.Get("env"),
		Verbose: verbose}, nil
}

// FormatDSN serialize a config to connect string
func (cfg *Config) FormatDSN() string {
	return fmt.Sprintf(`%s:%s@%s?env=%s&verbose=%s`, cfg.POPAccessID, cfg.POPAccessKey, cfg.POPURL, cfg.Env, strconv.FormatBool(cfg.Verbose))
}
