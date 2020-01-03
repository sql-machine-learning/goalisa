// Copyright 2020 The SQLFlow Authors. All rights reserved.
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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlisaConn_Exec(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("alisa", newConfigFromEnv(t).FormatDSN())
	a.NoError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sqlflow_alisa_test_table(c1 STRING);`)
	a.NoError(err)
}

func TestAlisaConn_Query(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("alisa", newConfigFromEnv(t).FormatDSN())
	a.NoError(err)

	rows, err := db.Query(`select "Alice" as name, 23.8 as age, 56000 as salary;`)
	a.NoError(err)
	defer rows.Close()

	columnNames, err := rows.Columns()
	a.NoError(err)
	a.Equal([]string{"name", "age", "salary"}, columnNames)

	columnTypes, err := rows.ColumnTypes()
	a.NoError(err)
	a.Equal(3, len(columnTypes))
	a.Equal("STRING", columnTypes[0].DatabaseTypeName())
	a.Equal("DOUBLE", columnTypes[1].DatabaseTypeName())
	a.Equal("BIGINT", columnTypes[2].DatabaseTypeName())

	rowContent := [][]string{}
	for rows.Next() {
		var c1, c2, c3 string
		err := rows.Scan(&c1, &c2, &c3)
		a.NoError(err)
		rowContent = append(rowContent, []string{c1, c2, c3})
	}
	a.Equal(1, len(rowContent))
	a.Equal("Alice", rowContent[0][0])
	a.Equal("23.8", rowContent[0][1])
	a.Equal("56000", rowContent[0][2])
}
