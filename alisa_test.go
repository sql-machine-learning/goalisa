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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var errSkipTesting = fmt.Errorf("skip test")

func TestCreateTask(t *testing.T) {
	if os.Getenv("POP_SECRET") == "" {
		t.Skip()
	}
	a := assert.New(t)
	ali := NewAlisaFromEnv()
	code := "SELECT 2;"
	taskID, _, err := ali.CreateTask(code)
	time.Sleep(time.Second * 2) // to avoid touching the flow-control
	a.NoError(err)
	a.NotEmpty(taskID)
	fmt.Println(ali.ReadLogs(taskID, 0))
}
