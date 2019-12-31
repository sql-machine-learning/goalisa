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
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type alisa struct {
	popURL    string
	popID     string
	popSecret string
	verbose   string
	envs      map[string]string
	pop       *popClient
}

// newAlisa init an Alisa client
func newAlisa(popURL, popID, popSecret, verbose, b64env string) (*alisa, error) {
	buf, err := base64.URLEncoding.DecodeString(b64env)
	if err != nil {
		return nil, err
	}
	var envs map[string]string
	if err := json.Unmarshal(buf, &envs); err != nil {
		return nil, err
	}
	return &alisa{
		popURL,
		popID,
		popSecret,
		verbose,
		envs,
		newPOP(-1),
	}, nil
}

// createTask returns a task id
func (ali *alisa) createTask(code string) (string, error) {
	params := baseParams(ali.popID)
	params["Action"] = "CreateAlisaTask"
	params["ExecCode"] = code

	params["SHOW_COLUMN_TYPE"] = "true" // display column type, for feature derivation.
	params["CustomerId"] = "sqlFlow"
	params["PluginName"] = "odps_sql"
	params["Exec"] = "/opt/taobao/tbdpapp/dwcommonwrapper/dwcommonwrapper.sh"
	params["UniqueKey"] = fmt.Sprintf("%d", time.Now().UnixNano())
	params["ExecTarget"] = ali.envs["ALISA_TASK_EXEC_TARGET"]
	envBuf, _ := json.Marshal(ali.envs)
	params["Envs"] = string(envBuf)
	rspBuf, err := ali.pop.request(params, ali.popURL, ali.popSecret)
	if err != nil {
		return "", err
	}
	// TODO(weiguo): return rspBuf["retValue"]["taskId"] in fact.
	return string(rspBuf), nil
}

// getStatus: returns the task status of taskID
func (ali *alisa) getStatus(taskID string) int {
	return 0
}

// completed: check if the status is completed
func (ali *alisa) completed(status int) bool {
	return true
}

// readLogs: reads task logs from `offset`
// return -1: read to the end
// return n(>0): keep reading with the offset `n` in the next time
func (ali *alisa) readLogs(taskID string, offset int) int {
	return 0
}

// readResults: reads the task results
func (ali *alisa) readResults(taskID string) {
	// TODO(weiguoz): define a result
}

// stop: stops the task
func (ali *alisa) stop(taskID string) bool {
	return true
}

func (ali *alisa) requetAndParseResult(params map[string]string) {
	return
}

func baseParams(popID string) map[string]string {
	gmt, _ := time.LoadLocation("GMT")
	uu, _ := uuid.NewUUID()
	return map[string]string{
		"Timestamp":        time.Now().In(gmt).Format(time.RFC3339),
		"AccessKeyId":      popID,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   strings.Replace(uu.String(), "-", "", -1),
		"Format":           "JSON",
		"Product":          "dataworks",
		"Version":          "2017-12-12",
	}
}