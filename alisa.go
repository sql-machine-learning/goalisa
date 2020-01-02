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

const maxLogNum = 1000

type alisa struct {
	popURL    string
	popID     string
	popSecret string
	verbose   bool
	envs      map[string]string
	pop       *popClient
}

// newAlisa init an Alisa client
func newAlisa(popURL, popID, popSecret, b64env string, verbose bool) (*alisa, error) {
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

// createTask returns a task id and it's status
func (ali *alisa) createTask(code string) (string, int, error) {
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

	res, err := ali.requetAndParseResponse(params)
	if err != nil {
		return "", -1, err
	}
	var val alisaTaskMeta
	if err = json.Unmarshal(*res, &val); err != nil {
		return "", -1, err
	}
	return val.TaskID, val.Status, nil
}

// getStatus: returns the task status of taskID
func (ali *alisa) getStatus(taskID string) (int, error) {
	params := baseParams(ali.popID)
	params["Action"] = "GetAlisaTask"
	params["AlisaTaskId"] = taskID

	res, err := ali.requetAndParseResponse(params)
	if err != nil {
		return -1, err
	}
	var val alisaTaskStatus
	if err = json.Unmarshal(*res, &val); err != nil {
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
	// the `maxLogs` used to deal with too many logs.
	end := false
	for i := 0; i < maxLogNum && !end; i++ {
		params := baseParams(ali.popID)
		params["Action"] = "GetAlisaTaskLog"
		params["AlisaTaskId"] = taskID
		params["Offset"] = fmt.Sprintf("%d", offset)

		res, err := ali.requetAndParseResponse(params)
		if err != nil {
			return offset, err
		}
		var log alisaTaskLog
		if err = json.Unmarshal(*res, &log); err != nil {
			return offset, err
		}
		rdLen, err := strconv.Atoi(log.ReadLen)
		if err != nil {
			return offset, err
		}
		if rdLen == 0 {
			return offset, nil
		}
		offset += rdLen
		end = log.End
		fmt.Printf(log.Content)
	}
	if end {
		return -1, nil
	}
	return offset, nil
}

func (ali *alisa) countResults(taskID string) (int, error) {
	params := baseParams(ali.popID)
	params["Action"] = "GetAlisaTaskResultCount"
	params["AlisaTaskId"] = taskID

	res, err := ali.requetAndParseResponse(params)
	if err != nil {
		return 0, err
	}
	var count string
	if err = json.Unmarshal(*res, &count); err != nil {
		return 0, err
	}
	return strconv.Atoi(count)
}

// readResults: reads the task results
func (ali *alisa) getResults(taskID string, batch int) (*alisaTaskResult, error) {
	if batch <= 0 {
		return nil, fmt.Errorf("batch shoud be lt 0")
	}
	nResults, err := ali.countResults(taskID)
	if err != nil {
		return nil, err
	}
	var taskRes alisaTaskResult
	for i := 0; i < nResults; i += batch {
		params := baseParams(ali.popID)
		params["Action"] = "GetAlisaTaskResult"
		params["AlisaTaskId"] = taskID
		params["Start"] = fmt.Sprintf("%d", i)
		params["Limit"] = fmt.Sprintf("%d", batch)
		res, err := ali.requetAndParseResponse(params)
		if err != nil {
			return nil, err
		}
		parseAlisaTaskResult(res, &taskRes)
	}
	return &taskRes, nil
}

// stop: stops the task.
// TODO(weiguz): need more tests
func (ali *alisa) stop(taskID string) (bool, error) {
	params := baseParams(ali.popID)
	params["Action"] = "StopAlisaTask"
	params["AlisaTaskId"] = taskID
	res, err := ali.requetAndParseResponse(params)
	if err != nil {
		return false, err
	}
	var ok bool
	if err = json.Unmarshal(*res, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (ali *alisa) requetAndParseResponse(params map[string]string) (*json.RawMessage, error) {
	rspBuf, err := ali.pop.request(params, ali.popURL, ali.popSecret)
	if err != nil {
		return nil, err
	}
	var aliRsp alisaResponse
	if err = json.Unmarshal(rspBuf, &aliRsp); err != nil {
		return nil, err
	}
	if aliRsp.Code != "0" {
		return nil, fmt.Errorf("bad result, response=%s", string(rspBuf))
	}
	return aliRsp.Value, nil
}

func parseAlisaTaskResult(from *json.RawMessage, to *alisaTaskResult) error {
	var rawResult alisaTaskRawResult
	if err := json.Unmarshal(*from, &rawResult); err != nil {
		return err
	}

	if len(to.columns) == 0 {
		bytHeader := []byte(rawResult.Header)
		var header []string
		if err := json.Unmarshal(bytHeader, &header); err != nil {
			return err
		}
		for _, h := range header {
			nt := strings.Split(h, "::")
			if len(nt) != 2 {
				return fmt.Errorf("bad header of alisa task result")
			}
			to.columns = append(to.columns, alisaTaskResultColumn{name: nt[0], typ: nt[1]})
		}
	}

	bytBody := []byte(rawResult.Body)
	var body [][]string
	if err := json.Unmarshal(bytBody, &body); err != nil {
		return err
	}
	for _, line := range body {
		to.body = append(to.body, line)
	}
	return nil
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
