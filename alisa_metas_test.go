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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlisaTaskMeta(t *testing.T) {
	a := assert.New(t)
	byt := []byte(`{
		"returnCode":"0",
		"returnValue":{"alisaTaskId":"G3_1","createTime":1577759588481,"logPath":"/path/to/log","localLogPath":"/20191231/path/to/your/local/","status":1},
		"requestId":"C3D87399-447E-4C53-B008-C1787C379D2E",
		"returnMessage":"ok",
		"returnErrorSolution":"keep going"
	}`)
	var rsp alisaResponse
	err := json.Unmarshal(byt, &rsp)
	a.NoError(err)
	a.Equal(rsp.Code, "0")

	var val alisaTaskMeta
	err = json.Unmarshal(*rsp.Value, &val)
	a.NoError(err)
	a.Equal(val.TaskID, "G3_1")
}

func TestAlisaTaskStatus(t *testing.T) {
	a := assert.New(t)
	byt := []byte(`{
		"returnCode":"0",
		"returnValue":{"nodeName":"alisa-test011167179122.na62sqa","alisaTaskId":"T3_0008271210","filePath":"/20191231/sqlFlow/13/27/00/ab09u281vsw9dmgh6kyimmvt/","nodeAddress":"11.167.179.122","status":2,"ossPath":"prod/sqlFlow/2019/12/31/6reg65oxvtng7hwxnyimyqs4"},
		"requestId":"0b9e82cd15777700362005616d1603",
		"returnMessage":"ok",
		"returnErrorSolution":"keep going"
	}`)
	var rsp alisaResponse
	err := json.Unmarshal(byt, &rsp)
	a.NoError(err)
	a.Equal(rsp.Code, "0")

	var val alisaTaskStatus
	err = json.Unmarshal(*rsp.Value, &val)
	a.NoError(err)
	a.Equal(val.Status, 2)
}

func TestAlisaTaskLog(t *testing.T) {
	a := assert.New(t)
	byt := []byte(`{
		"returnCode":"0",
		"returnValue":{"readLength":"1","logMsg":"R","isEnd":true},
		"requestId":"0b9e8d8215777702027777560d135b",
		"returnMessage":"ok",
		"returnErrorSolution":"keep going"
	}`)
	var rsp alisaResponse
	err := json.Unmarshal(byt, &rsp)
	a.NoError(err)
	a.Equal(rsp.Code, "0")

	var val alisaTaskLog
	err = json.Unmarshal(*rsp.Value, &val)
	a.NoError(err)
	a.Equal(val.End, true)
}

func TestAlisaTaskResultCount(t *testing.T) {
	a := assert.New(t)
	byt := []byte(`{
		"returnCode":"0",
		"returnValue":"1",
		"requestId":"0b9e8d8215777650591238870d4f31",
		"returnMessage":"ok",
		"returnErrorSolution":"keep going"
	}`)
	var rsp alisaResponse
	err := json.Unmarshal(byt, &rsp)
	a.NoError(err)
	a.Equal(rsp.Code, "0")

	var count string
	err = json.Unmarshal(*rsp.Value, &count)
	a.NoError(err)
	a.Equal(count, "1")
}

func TestAlisaTaskResult(t *testing.T) {
	a := assert.New(t)
	byt := []byte(`{
		"returnCode":"0",
		"returnValue":{"dataHeader":"[\"name::string\",\"age::double\",\"salary::bigint\"]","resultMsg":"[[\"3m^2\",\"23.8\",\"56000\"]]"},
		"requestId":"0b9e8d8215777683831622186d4f31",
		"returnMessage":"ok",
		"returnErrorSolution":"keep going"
	}`)
	var rsp alisaResponse
	err := json.Unmarshal(byt, &rsp)
	a.NoError(err)
	a.Equal(rsp.Code, "0")

	var val alisaTaskRawResult
	err = json.Unmarshal(*rsp.Value, &val)
	a.NoError(err)

	bytHeader := []byte(val.Header)
	var header []string
	err = json.Unmarshal(bytHeader, &header)
	a.NoError(err)
	a.Equal(len(header), 3)

	bytBody := []byte(val.Body)
	var body [][]string
	err = json.Unmarshal(bytBody, &body)
	a.NoError(err)
	a.Equal(len(body), 1)
	a.Equal(len(body[0]), len(header))
}
