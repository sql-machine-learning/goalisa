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
	"database/sql"
	"database/sql/driver"
)

// register driver
func init() {
	sql.Register("alisa", &Driver{})
}

// Driver implements database/sql/driver.Driver
type Driver struct{}

// Open returns a connection
func (dr Driver) Open(dsn string) (driver.Conn, error) {
	config, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	alisa := NewAlisa(config)
	if err != nil {
		return nil, err
	}

	return &alisaConn{alisa}, nil
}
