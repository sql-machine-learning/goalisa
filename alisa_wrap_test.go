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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryAlisaTask(t *testing.T) {
	a := assert.New(t)
	ali := newAlisaFromEnv(t)
	cmd := "select \"Alice\" as name, 23.8 as age, 56000 as salary;"
	res, err := ali.query(cmd)
	a.NoError(err)
	// schema, header
	a.Equal(len(res.columns), 3)
	a.Equal(res.columns[0].name, "name")
	a.Equal(res.columns[0].typ, "string")
	a.Equal(res.columns[1].name, "age")
	a.Equal(res.columns[1].typ, "double")
	a.Equal(res.columns[2].name, "salary")
	a.Equal(res.columns[2].typ, "bigint")
	// body
	a.Equal(len(res.body), 1)
	a.Equal(len(res.body[0]), 3)
	a.Equal(res.body[0][0], "Alice")
	a.Equal(res.body[0][1], "23.8")
	a.Equal(res.body[0][2], "56000")
}

func TestExecAlisaTask(t *testing.T) {
	a := assert.New(t)
	ali := newAlisaFromEnv(t)
	cmd := "CREATE TABLE IF NOT EXISTS sqlflow_alisa_test_table(c1 STRING);"
	err := ali.exec(cmd)
	a.NoError(err)
}
