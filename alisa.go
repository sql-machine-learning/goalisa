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
		&popClient{},
	}, nil
}

// return task id
func (ali *alisa) createTask(code string) (string, error) {
	params := ali.baseParams()
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
	rsp, err := ali.pop.request(params, ali.popURL, ali.popSecret)
	if err != nil {
		return "", err
	}
	// TODO(weiguo): return rsp["retValue"]["taskId"] in fact.
	return rsp, nil
}

func (ali *alisa) baseParams() map[string]string {
	gmt, _ := time.LoadLocation("GMT")
	uu, _ := uuid.NewUUID()
	return map[string]string{
		"Timestamp":        time.Now().In(gmt).Format(time.RFC3339),
		"AccessKeyId":      ali.popID,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   strings.Replace(uu.String(), "-", "", -1),
		"Format":           "JSON",
		"Product":          "dataworks",
		"Version":          "2017-12-12",
	}
}
