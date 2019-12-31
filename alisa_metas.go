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

import "encoding/json"

type alisaTaskResultColumn struct {
	name string
	typ  string
}

type alisaTaskResult struct {
	columns []alisaTaskResultColumn
	body    [][]string
}

type alisaResponse struct {
	RequestID     string           `json:"requestId"`
	Code          string           `json:"returnCode"`
	Value         *json.RawMessage `json:"returnValue"`
	Message       string           `json:"returnMessage"`
	ErrorSolution string           `json:"returnErrorSolution"`
}

// alisaResponse.Value
type alisaTaskMeta struct {
	TaskID       string `json:"alisaTaskId"`
	CreateTime   int64  `json:"createTime"`
	Status       int    `json:"status"`
	LogPath      string `json:"logPath"`
	LocalLogPath string `json:"localLogPath"`
}

// alisaResponse.Value
type alisaTaskStatus struct {
	TaskID      string `json:"alisaTaskId"`
	Status      int    `json:"status"`
	NodeName    string `json:"nodeName"`
	FilePath    string `json:"filePath"`
	NodeAddress string `json:"nodeAddress"`
	OssPath     string `json:"ossPath"`
}

// alisaResponse.Value
type alisaTaskLog struct {
	ReadLen string `json:"readLength"`
	Content string `json:"logMsg"`
	End     bool   `json:"isEnd"`
}

// alisaResponse.Value
// taskResultCount: single string

// alisaResponse.Value
type alisaTaskRawResult struct {
	Header string `json:"dataHeader"`
	Body   string `json:"resultMsg"`
}
