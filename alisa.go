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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// alisaTask*: task status returned by getStatus()
	alisaTaskWaiting    = 1
	alisaTaskRunning    = 2
	alisaTaskCompleted  = 3
	alisaTaskError      = 4
	alisaTaskFailover   = 5
	alisaTaskKilled     = 6
	alisaTaskRerun      = 8
	alisaTaskExpired    = 9
	alisaTaskAlisaRerun = 10
	alisaTaskAllocate   = 11

	// used to deal with too many logs.
	maxLogNum = 2000
)

type Alisa struct {
	*Config
	pop *popClient
}

// NewAlisaFromEnv initializes Alisa client
func NewAlisaFromEnv() *Alisa {
	cfg := newConfigFromEnv()
	return &Alisa{cfg, newPOP(-1)}
}

// NewAlisa init an Alisa client
func NewAlisa(cfg *Config) *Alisa {
	return &Alisa{cfg, newPOP(-1)}
}

func newConfigFromEnv() *Config {
	popURL := os.Getenv("POP_URL")
	popID := os.Getenv("POP_ID")
	popSecret := os.Getenv("POP_SECRET")
	verbose := len(os.Getenv("VERBOSE")) > 0
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
	return &Config{POPAccessID: popID, POPAccessSecret: popSecret, POPURL: popURL, Verbose: verbose, Env: envs}
}

// CreateTask returns a task id and it's status
func (ali *Alisa) CreateTask(code string) (string, int, error) {
	params := baseParams(ali.POPAccessID)
	params["ExecCode"] = code

	params["SHOW_COLUMN_TYPE"] = "true" // display column type, for feature derivation.
	params["CustomerId"] = "sqlFlow"
	params["PluginName"] = "odps_sql"
	params["Exec"] = "/opt/taobao/tbdpapp/dwcommonwrapper/dwcommonwrapper.sh"
	params["UniqueKey"] = fmt.Sprintf("%d", time.Now().UnixNano())
	params["ExecTarget"] = ali.Env["ALISA_TASK_EXEC_TARGET"]
	envBuf, _ := json.Marshal(ali.Env)
	params["Envs"] = string(envBuf)
	res, err := ali.requetAndParseResponse("CreateAlisaTask", params)
	if err != nil {
		return "", -1, err
	}
	var val alisaTaskMeta
	if err = json.Unmarshal(*res, &val); err != nil {
		return "", -1, err
	}
	return val.TaskID, val.Status, nil
}

// GetStatus returns the task status of taskID
func (ali *Alisa) GetStatus(taskID string) (int, error) {
	params := baseParams(ali.POPAccessID)
	params["AlisaTaskId"] = taskID
	res, err := ali.requetAndParseResponse("GetAlisaTask", params)
	if err != nil {
		return -1, err
	}
	var val alisaTaskStatus
	if err = json.Unmarshal(*res, &val); err != nil {
		return -1, err
	}
	// alisaTask*
	return val.Status, nil
}

// Completed checks if the status is completed
func (ali *Alisa) Completed(status int) bool {
	return status == alisaTaskCompleted || status == alisaTaskError || status == alisaTaskKilled || status == alisaTaskRerun || status == alisaTaskExpired
}

// ReadLogs reads task logs from `offset`
// return -1: read to the end
// return n(>0): keep reading with the offset `n` in the next time
func (ali *Alisa) ReadLogs(taskID string, offset int) (int, error) {
	end := false
	for i := 0; i < maxLogNum && !end; i++ {
		params := baseParams(ali.POPAccessID)
		params["AlisaTaskId"] = taskID
		params["Offset"] = fmt.Sprintf("%d", offset)
		res, err := ali.requetAndParseResponse("GetAlisaTaskLog", params)
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
		fmt.Print(log.Content)
	}
	if end {
		return -1, nil
	}
	return offset, nil
}

func (ali *Alisa) countResults(taskID string) (int, error) {
	params := baseParams(ali.POPAccessID)
	params["AlisaTaskId"] = taskID
	res, err := ali.requetAndParseResponse("GetAlisaTaskResultCount", params)
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
func (ali *Alisa) getResults(taskID string, batch int) (*alisaTaskResult, error) {
	if batch <= 0 {
		return nil, fmt.Errorf("batch shoud be lt 0")
	}
	nResults, err := ali.countResults(taskID)
	if err != nil {
		return nil, err
	}
	var taskRes alisaTaskResult
	for i := 0; i < nResults; i += batch {
		params := baseParams(ali.POPAccessID)
		params["AlisaTaskId"] = taskID
		params["Start"] = fmt.Sprintf("%d", i)
		params["Limit"] = fmt.Sprintf("%d", batch)
		res, err := ali.requetAndParseResponse("GetAlisaTaskResult", params)
		if err != nil {
			return nil, err
		}
		parseAlisaTaskResult(res, &taskRes)
	}
	return &taskRes, nil
}

// stop: stops the task.
// TODO(weiguz): need more tests
func (ali *Alisa) stop(taskID string) (bool, error) {
	params := baseParams(ali.POPAccessID)
	params["AlisaTaskId"] = taskID
	res, err := ali.requetAndParseResponse("StopAlisaTask", params)
	if err != nil {
		return false, err
	}
	var ok bool
	if err = json.Unmarshal(*res, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (ali *Alisa) requetAndParseResponse(action string, params map[string]string) (*json.RawMessage, error) {
	params["Action"] = action
	rspBuf, err := ali.pop.request(params, ali.POPURL, ali.POPAccessSecret)
	if err != nil {
		return nil, fmt.Errorf("%s got an error: %v", action, err)
	}
	var aliRsp alisaResponse
	if err = json.Unmarshal(rspBuf, &aliRsp); err != nil {
		return nil, fmt.Errorf("%s got an error: %v", action, err)
	}
	if aliRsp.Code != "0" {
		return nil, fmt.Errorf("%s got a bad result, response=%s", action, string(rspBuf))
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
