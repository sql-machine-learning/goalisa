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
	"reflect"
)

var builtinString = reflect.TypeOf(string(""))

type alisaRows struct {
	// TODO(weiguoz)
}

func (ar *alisaRows) Close() error {
	// TODO(weiguoz)
	return nil
}

func (ar *alisaRows) Columns() []string {
	// TODO(weiguoz)
	return nil
}

// Notice: `\N` denotes nil, even outher types
func (ar *alisaRows) Next(dst []driver.Value) error {
	// TODO(weiguoz)
	return nil
}

func (ar *alisaRows) ColumnTypeScanType(i int) reflect.Type {
	// TODO(weiguoz)
	return builtinString
}

func (ar *alisaRows) ColumnTypeDatabaseTypeName(i int) string {
	// TODO(weiguoz)
	return ""
}
