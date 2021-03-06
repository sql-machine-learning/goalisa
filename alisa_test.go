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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var errSkipTesting = fmt.Errorf("skip test")

func newAlisaFromEnv(t *testing.T) *Alisa {
	cfg := newConfigFromEnv(t)
	return New(cfg)
}

func TestCreateSQLTask(t *testing.T) {
	a := assert.New(t)
	ali := newAlisaFromEnv(t)
	code := "SELECT 2;"
	taskID, _, err := ali.createSQLTask(code)
	time.Sleep(time.Second * 2) // to avoid touching the flow-control
	a.NoError(err)
	a.NotEmpty(taskID)
}

func TestCreatePyODPSTask(t *testing.T) {
	a := assert.New(t)
	ali := newAlisaFromEnv(t)
	code := `import argparse

if __name__ == "__main__":
    input_table_name = args['input_table']
    output_table_name = args['output_table']
    print(input_table_name)
    print(output_table_name)

    input_table = o.get_table(input_table_name)
    print(input_table.schema)
    output_table = o.create_table(output_table_name, input_table.schema)
	`
	args := "input_table=table_1 output_table=table_2"
	taskID, _, err := ali.createPyODPSTask(code, args)
	time.Sleep(time.Second * 2) // to avoid touching the flow-control
	a.NoError(err)
	a.NotEmpty(taskID)
}
