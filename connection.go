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
	"database/sql/driver"
	"net/http"
)

type alisaConn struct {
	*http.Client
	*Config
}

// Begin, unimplemented
func (ac *alisaConn) Begin() (driver.Tx, error) {
	panic("Not implemented")
}

// Prepare, unimplemented
func (ac *alisaConn) Prepare(query string) (driver.Stmt, error) {
	panic("Not implemented")
}

// Close connection
func (ac *alisaConn) Close() error {
	// TODO(weiguoz)
	return nil
}

// Exec implements database/sql/driver.Execer.
// Note: result is always nil
func (ac *alisaConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	// TODO(weiguoz), return alisaResult
	return &alisaResult{-1, -1}, nil
}

// Query implements database/sql/driver.Queryer
func (ac *alisaConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	// TODO(weiguo), returns rows
	return &alisaRows{}, nil
}
