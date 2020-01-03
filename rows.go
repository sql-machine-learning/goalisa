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
	"io"
	"reflect"
	"strings"
)

var builtinString = reflect.TypeOf(string(""))

type alisaRows struct {
	rowIdx int
	result *alisaTaskResult
}

// Close closes the rows iterator.
func (ar *alisaRows) Close() error {
	ar.rowIdx = -1
	ar.result = nil
	return nil
}

// Columns returns the names of the columns.
func (ar *alisaRows) Columns() []string {
	columnNames := []string{}
	for _, c := range ar.result.columns {
		columnNames = append(columnNames, c.name)
	}
	return columnNames
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// NOTE(weiguoz): `\N` denotes nil, even outher types
func (ar *alisaRows) Next(dst []driver.Value) error {
	if ar.rowIdx >= len(ar.result.body) {
		return io.EOF
	}

	// Fill in dest with one single row data.
	for colIndex, value := range ar.result.body[ar.rowIdx] {
		dst[colIndex] = value
	}

	ar.rowIdx++
	return nil
}

// RowsColumnTypeScanType always gives string type
func (ar *alisaRows) ColumnTypeScanType(i int) reflect.Type {
	return builtinString
}

// RowsColumnTypeDatabaseTypeName returns the database system type name in uppercase.
func (ar *alisaRows) ColumnTypeDatabaseTypeName(i int) string {
	return strings.ToUpper(ar.result.columns[i].typ)
}
