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
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	a := assert.New(t)
	popURL := os.Getenv("POP_URL")
	popID := os.Getenv("POP_ID")
	popSecret := os.Getenv("POP_SECRET")
	verbose := "TBD"
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
	if len(envs["SKYNET_ACCESSKEY"]) == 0 {
		t.Skip()
	}
	rawEnv, err := json.Marshal(envs)
	a.NoError(err)
	b64Envs := base64.StdEncoding.EncodeToString(rawEnv)

	alisa, err := newAlisa(popURL, popID, popSecret, verbose, b64Envs)
	a.NoError(err)
	code := "SELECT 2;"
	taskID, err := alisa.createTask(code)
	a.NoError(err)
	a.NotEmpty(taskID)
}
