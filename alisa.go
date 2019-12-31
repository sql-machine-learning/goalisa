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
	"strconv"
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
	var aliRsp alisaResponse
	if err = json.Unmarshal(rspBuf, &aliRsp); err != nil {
		return "", err
	}
	if aliRsp.Code != "0" {
		return "", fmt.Errorf("bad result, code=%s, message=%s", aliRsp.Code, aliRsp.Message)
	}
	var val taskMeta
	if err = json.Unmarshal(*aliRsp.Value, &val); err != nil {
		return "", err
	}
	return val.TaskID, nil
}

// getStatus: returns the task status of taskID
func (ali *alisa) getStatus(taskID string) (int, error) {
	params := baseParams(ali.popID)
	params["Action"] = "GetAlisaTask"
	params["AlisaTaskId"] = taskID

	rspBuf, err := ali.pop.request(params, ali.popURL, ali.popSecret)
	if err != nil {
		return -1, err
	}
	var aliRsp alisaResponse
	if err = json.Unmarshal(rspBuf, &aliRsp); err != nil {
		return -1, err
	}
	if aliRsp.Code != "0" {
		return -1, fmt.Errorf("bad result, code=%s, message=%s", aliRsp.Code, aliRsp.Message)
	}
	var val taskStatus
	if err = json.Unmarshal(*aliRsp.Value, &val); err != nil {
		return -1, err
	}
	return val.Status, nil
}

// completed: check if the status is completed
func (ali *alisa) completed(status int) bool {
	return status == 3 || status == 4 || status == 6 || status == 8 || status == 9
}

// readLogs: reads task logs from `offset`
// return -1: read to the end
// return n(>0): keep reading with the offset `n` in the next time
func (ali *alisa) readLogs(taskID string, offset int) (int, error) {
	// we don't trust the server returned `end`, so the `maxLogs` used to deal with too many logs.
	end := false
	for maxLogs := 10000; !end && maxLogs > 0; maxLogs-- {
		params := baseParams(ali.popID)
		params["Action"] = "GetAlisaTaskLog"
		params["AlisaTaskId"] = taskID
		params["Offset"] = fmt.Sprintf("%d", offset)
		rspBuf, err := ali.pop.request(params, ali.popURL, ali.popSecret)
		if err != nil {
			return 0, err
		}
		var aliRsp alisaResponse
		if err = json.Unmarshal(rspBuf, &aliRsp); err != nil {
			return 0, err
		}
		if aliRsp.Code != "0" {
			return 0, fmt.Errorf("bad result, code=%s, message=%s", aliRsp.Code, aliRsp.Message)
		}
		var val taskLog
		if err = json.Unmarshal(*aliRsp.Value, &val); err != nil {
			return 0, err
		}
		rdLen, err := strconv.Atoi(val.ReadLen)
		if err != nil {
			return 0, nil
		}
		if rdLen == 0 {
			return offset, nil
		}
		offset += rdLen
		end = val.End
		fmt.Printf(val.Content)
	}
	if end {
		return -1, nil
	}
	return offset, nil
}

// readResults: reads the task results
func (ali *alisa) getResults(taskID string) error {
	// TODO(weiguoz): define a result
	return nil
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
